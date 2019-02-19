package grest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

//ServeReadCloser returns a HTTP response with Content coming from a io.ReadCloser that is closed after Read() returns io.EOF
//If getReader() returns an error it will result in a panic
//Try to create the reader inside the getReader func to avoid too soon/unnecessary memory allocation
func ServeReadCloser(getReader func() (io.ReadCloser, error)) WebPart {
	return func(u WebUnit) *WebUnit {
		r, err := getReader()
		if err != nil {
			u.Panic(err)
			u.Writer.WriteHeader(http.StatusInternalServerError)
			u.Writer.Write([]byte(err.Error()))
			return &u
		}
		defer r.Close()
		status := u.GetStatus()
		if status == 0 {
			status = http.StatusOK
		}
		u.Writer.WriteHeader(status)
		_, err = io.Copy(u.Writer, r)
		if err != nil {
			u.Writer.WriteHeader(http.StatusInternalServerError)
			u.Writer.Write([]byte(err.Error()))
		}

		return &u
	}
}

//ServeReadCloser returns a HTTP response with Content coming from a io.ReadCloser that is closed after Read() returns io.EOF
//If getReader() returns an error it will result in a panic
func (w WebPart) ServeReadCloser(getReader func() (io.ReadCloser, error)) WebPart {
	return Compose(w, ServeReadCloser(getReader))
}

// ServeBytes responses with the given bytes
func ServeBytes(data []byte) WebPart {
	return ServeReadCloser(func() (io.ReadCloser, error) { return MakeClosable(bytes.NewReader(data), nil), nil })
}

// ServeBytes responses with the given bytes
func (w WebPart) ServeBytes(data []byte) WebPart {
	return Compose(w, ServeBytes(data))
}

//ServeString serves the given string as response (convinience wrapper for ServeBytes)
func ServeString(s string) WebPart {
	return ServeBytes([]byte(s))
}

//ServeString serves the given string as response (convinience wrapper for ServeBytes)
func (w WebPart) ServeString(s string) WebPart {
	return w.ServeBytes([]byte(s))
}

//ServeJSON responses with a JSON object as bytes
func ServeJSON(obj interface{}) WebPart {
	return ServeReadCloser(func() (io.ReadCloser, error) {
		data, err := json.Marshal(obj)
		if err != nil {
			return nil, err
		}
		return MakeClosable(bytes.NewReader(data), nil), nil
	})
}

//ServeJSON responses with a JSON object as bytes
func (w WebPart) ServeJSON(obj interface{}) WebPart {
	return Compose(w, ServeJSON(obj))
}

//ServeExtrasAsJSON is a convinience call that converts the current Extras into a JSON object and returns it with the current status (default 200 OK)
func ServeExtrasAsJSON() WebPart {
	return func(u WebUnit) *WebUnit {
		return ServeJSON(u.Extras())(u)
	}
}

//ServeExtrasAsJSON is a convinience call that converts the current Extras into a JSON object and returns it with the current status (default 200 OK)
func (w WebPart) ServeExtrasAsJSON() WebPart {
	return Compose(w, ServeExtrasAsJSON())
}

//=== Helpers =====================================================================================
type closer struct {
	reader io.Reader
	close  func() error
}

func (c closer) Read(buffer []byte) (int, error) {
	return c.reader.Read(buffer)
}

func (c closer) Close() error {
	if c.close != nil {
		return c.close()
	}
	return nil
}

//MakeClosable turns a io.Reader into io.ReadCloser with optional closeAction
func MakeClosable(r io.Reader, closeAction func() error) io.ReadCloser {
	return &closer{r, closeAction}
}
