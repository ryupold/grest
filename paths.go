package grest

import (
	"regexp"
	"strconv"
	"strings"
)

// Prefix filters paths that dont start with 'prefix'
func Prefix(prefix string) WebPart {
	return func(u WebUnit) *WebUnit {
		if strings.HasPrefix(Clean(u.Request.URL.Path), prefix) {
			return &u
		}
		return nil
	}
}

// Prefix (composing) filters paths that dont start with 'prefix'
func (w WebPart) Prefix(prefix string) WebPart {
	return Compose(w, Prefix(prefix))
}

// PrefixDirty filters paths that dont start with 'prefix' (doesn't clean path)
func PrefixDirty(prefix string) WebPart {
	return func(u WebUnit) *WebUnit {
		if strings.HasPrefix(u.Request.URL.Path, prefix) {
			return &u
		}
		return nil
	}
}

// PrefixDirty (composing) filters paths that dont start with 'prefix' (doesn't clean path)
func (w WebPart) PrefixDirty(prefix string) WebPart {
	return Compose(w, PrefixDirty(prefix))
}

// Path matches exact path
func Path(path string) WebPart {
	return func(u WebUnit) *WebUnit {
		if Clean(u.Request.URL.Path) == path {
			return &u
		}
		return nil
	}
}

// Path (composing) matches exact path
func (w WebPart) Path(path string) WebPart {
	return Compose(w, Path(path))
}

// TypedPath parses an URL and says yes if it has all typed parameters as part of its part.
// e.g.: http://test.de/test/%s/%d is a pattern that expects the url to have a string as second part of the path
// and an integer as third part
// Allowed types:
// %s => string (any),
// %d => int,
// %f => float64,
// %t => bool
func TypedPath(pattern string, do func(WebUnit, []interface{}) *WebUnit) WebPart {
	return func(u WebUnit) *WebUnit {
		pParts := strings.Split(strings.Trim(pattern, "/"), "/")
		uParts := strings.Split(strings.Trim(u.Request.URL.Path, "/"), "/")

		if len(pParts) == len(uParts) {
			var values []interface{}
			for i, p := range pParts {
				switch p {
				//string
				case "%s":
					values = append(values, uParts[i])
				//int
				case "%d":
					v, e := strconv.Atoi(uParts[i])
					if e != nil {
						return nil
					}
					values = append(values, v)
				//float
				case "%f":
					v, e := strconv.ParseFloat(uParts[i], 64)
					if e != nil {
						return nil
					}
					values = append(values, v)
				//bool
				case "%t":
					v, e := strconv.ParseBool(uParts[i])
					if e != nil {
						return nil
					}
					values = append(values, v)

				//constant strings, must match exactly
				default:
					if pParts[i] != uParts[i] {
						return nil
					}
				}
			}

			return do(u, values)
		}
		return nil
	}
}

// TypedPath parses an URL and says yes if it has all typed parameters as part of its part.
// e.g.: http://test.de/test/%s/%d is a pattern that expects the url to have a string as second part of the path
// and an integer as third part
// Allowed types:
// %s => string (any),
// %d => int,
// %f => float64,
// %t => bool
func (w WebPart) TypedPath(pattern string, do func(WebUnit, []interface{}) *WebUnit) WebPart {
	return Compose(w, TypedPath(pattern, do))
}

// RegexPath matches path by regular expression
// e.g.: ^/[a-z]+[0-9]+$ matches http://website.de/test1
// if no match
func RegexPath(pattern string) WebPart {
	return func(u WebUnit) *WebUnit {
		path := u.Request.URL.Path
		if path != "/" {
			path = strings.TrimSuffix(path, "/")
		}
		if m, err := regexp.MatchString(pattern, path); m {
			return &u
		} else if err != nil {
			u.Panic(err)
			return &u
		}

		return nil
	}
}

// RegexPath matches path by regular expression
func (w WebPart) RegexPath(pattern string) WebPart {
	return Compose(w, RegexPath(pattern))
}

// Trim trims trailing slashes
func Trim() WebPart {
	return func(u WebUnit) *WebUnit {
		path := u.Request.URL.Path
		path = strings.Replace(strings.TrimSpace(path), "//", "/", -1)
		if path != "/" {
			path = strings.TrimSuffix(path, "/")
		}
		u.Request.URL.Path = path
		return &u
	}
}

// Trim (composing) trims trailing slashes
func (w WebPart) Trim() WebPart {
	return Compose(w, Trim())
}

// Clean makes path lowercase, removes // and the trailing /
func Clean(path string) string {
	path = strings.Replace(strings.ToLower(strings.TrimSpace(path)), "//", "/", -1)
	if path != "/" {
		path = strings.TrimSuffix(path, "/")
	}
	return path
}
