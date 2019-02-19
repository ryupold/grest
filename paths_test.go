package grest

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

type testWriter struct {
	HeaderValues map[string][]string
	Body         []byte
}

func (t testWriter) Header() http.Header {
	return t.HeaderValues
}

func (t testWriter) Write(data []byte) (int, error) {
	i := 0
	for _, d := range data {
		t.Body = append(t.Body, d)
		i++
	}

	return i, nil
}

func (t testWriter) WriteHeader(statusCode int) {
	t.Body = append(t.Body, []byte(fmt.Sprint(statusCode))...)
}

func getTestContext(useURL string) WebUnit {
	u, err := url.ParseRequestURI(useURL)
	try(err)
	return WebUnit{testWriter{HeaderValues: make(map[string][]string)}, &http.Request{URL: u}, context.TODO()}
}

func TestPath(t *testing.T) {
	cases := []struct {
		test string
		url  string
		pass bool
	}{
		{"/test", "/test", true},
		{"/test", "/test/", true},
		{"/test", "/TesT", true},
		{"/test", "/test/more", false},
		{"/test/more", "/test/more", true},
		{"/test/more", "/test//more", true},
	}

	for _, c := range cases {
		sut := Path(c.test)
		result := sut(getTestContext("http://text.de" + c.url))
		if (result == nil && c.pass) || (result != nil && !c.pass) {
			t.Errorf(`Path("%s") on URL=%s%s should be %v`, c.test, "http://text.de", c.url, c.pass)
		}
	}
}

func TestPrefix(t *testing.T) {
	cases := []struct {
		test string
		url  string
		pass bool
	}{
		{"/test", "/test", true},
		{"/test", "/test/", true},
		{"/test", "/TesT", true},
		{"/test", "/test/more", true},
		{"/test/more", "/test/more", true},
		{"/test/more", "/test//more", true},
		{"/test/more", "//more", false},
		{"/test", "//more", false},
	}

	for _, c := range cases {
		sut := Prefix(c.test)
		result := sut(getTestContext("http://text.de" + c.url))
		if (result == nil && c.pass) || (result != nil && !c.pass) {
			t.Errorf(`Prefix("%s") on URL=%s%s should be %v`, c.test, "http://text.de", c.url, c.pass)
		}
	}
}

func TestTypedPath(t *testing.T) {
	cases := []struct {
		pattern   string
		url       string
		values    []interface{}
		hasResult bool
	}{
		{"/test/%s/%d/%f/%t", "/test/string/5/5.5/false", []interface{}{"string", 5, 5.5, false}, true},
		{"/other/%s/%d/%f/%t", "/test/string/5/5.5/false", []interface{}{"string", 5, 5.5, false}, false},
		{"/test/%t", "/test/true", []interface{}{true}, true},
		{"/test/%t", "/test/string", []interface{}{"string"}, false},
		{"/test/%d", "/test/5.5", []interface{}{5.5}, false},
		{"/test/%f", "/test/5.5", []interface{}{5.5}, true},
		{"/test/%d", "/test/12345", []interface{}{12345}, true},
		{"/test/%s/%d/%f/%t", "/test/string/5/5.5", []interface{}{"string", 5, 5.5}, false},
		{"/test/%s/%d/%f/", "/test/string/5/5.5", []interface{}{"string", 5, 5.5}, true},

		{"/test/%f", "/test/5.5", []interface{}{5}, false},
		{"/test/%t/%s", "/test/true/asdasd", []interface{}{true, "asdasd"}, true},
		{"/test/%t/%s/%d", "/test/true/asdasd/5", []interface{}{true, "asdasd"}, false},
		{"/test/%t/%s/%d", "/test/true/asdasd/5", []interface{}{true, "asdasd", 5.5}, false},
		{"/test/%t/%s/%d", "/test/true/asdasd/5", []interface{}{true, "asdasd", 5}, true},
	}

	for _, c := range cases {
		sut := TypedPath(
			c.pattern,
			func(unit WebUnit, v []interface{}) *WebUnit {
				if len(v) == len(c.values) {
					for i := range v {
						if v[i] != c.values[i] {
							if c.hasResult {
								t.Errorf("Value %d differs: %v != %v\n", i, v[i], c.values[i])
							}
							return nil
						}
					}
					return &unit
				}
				if c.hasResult {
					t.Errorf("Parameter count differs: %v != %v\n", len(v), len(c.values))
				}
				return nil
			})
		result := sut(getTestContext("http://text.de" + c.url))
		if (result == nil && c.hasResult) || (result != nil && !c.hasResult) {
			t.Errorf(`TypedPath("%s") on URL=%s%s should be %v`, c.pattern, "http://text.de", c.url, c.hasResult)
		}
	}
}

func TestRequiredQuery(t *testing.T) {
	cases := []struct {
		keys      []string
		url       string
		hasResult bool
	}{
		{[]string{"test1"}, "?test1=5", true},
		{[]string{"test1"}, "?test1=5&test2=asd", true},
		{[]string{"test1", "test2"}, "?test1=5&test2=asd", true},
		{[]string{"test1", "test3"}, "?test1=5&test2=asd", false},
		{[]string{"test"}, "test=missingquestionmark", false},
		{[]string{"5"}, "?5=canAlsoBeAnIntApperantly", true},
		{[]string{"a", "b"}, "?a&b", false},
		{[]string{"a", "b"}, "?a=b&=5", false},
		{[]string{"a", "b"}, "?a=0&b=1", true},
	}

	for _, c := range cases {
		sut := Query(c.keys...)
		result := sut(getTestContext("http://text.de" + c.url))
		if (result == nil && c.hasResult) || (result != nil && !c.hasResult) {
			t.Errorf(`RequiredQuery("%v") on URL=%s%s should be %v`, c.keys, "http://text.de", c.url, c.hasResult)
		}
	}
}

type TestType struct {
	A string
	B string
}

func TestRegexPath(t *testing.T) {
	cases := []struct {
		test   string
		url    string
		pass   bool
		onFail WebPart
	}{
		{"^/test$", "/test", true, nil},
		{"^/test$", "/test/", true, nil},
		{"^/[a-z]+$", "/test", true, nil},
		{"^/[a-z]+[0-9]+$", "/test", false, nil},
		{"^/[a-z]+[0-9]{1,3}$", "/test514", true, nil},
		{"^/[a-z]+[0-9]{1,3}$", "/test514a", false, nil},
		{"^/[a-z]+[0-9]{2,}$", "/test514", true, nil},
		{"^/[a-z]+[0-9]{2,}$", "/asd", false, nil},
		{"^/[a-z]+[0-9]{2,}$", "/asd", true, BadRequest().ServeJSON(TestType{})},
	}

	for _, c := range cases {
		sut := RegexPath(c.test)
		result := sut(getTestContext("http://text.de" + c.url))
		if (result == nil && c.pass) || (result != nil && !c.pass) {
			if c.onFail == nil || (result != nil && c.onFail(*result) == nil) {
				t.Errorf(`RegexPath("%s") on URL=%s%s should be %v`, c.test, "http://text.de", c.url, c.pass)
			}
		}
	}
}
