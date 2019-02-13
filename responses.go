package grest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// ResponseR returns a HTTP response with 'status' and data as Reader
func ResponseR(status int, reader func() io.ReadCloser) WebPart {
	return func(u WebUnit) *WebUnit {
		u.Writer.WriteHeader(status)
		r := reader()
		defer r.Close()
		_, err := io.Copy(u.Writer, r)
		Try(err)
		return &u
	}
}

// Response returns a HTTP response with 'status' and 'data'
func Response(status int, data []byte) WebPart {
	return ResponseR(status, func() io.ReadCloser { return MakeClosable(bytes.NewReader(data), nil) })
}

// Response (composing) returns a HTTP response with 'status' and 'data'
func (w WebPart) Response(status int, data []byte) WebPart {
	return Compose(w, Response(status, data))
}

// ResponseJ returns a HTTP response with status and a json object as data
func ResponseJ(status int, v interface{}) WebPart {
	data, err := json.Marshal(v)
	Try(err)
	return ContentType(ContentTypeJSON).Response(status, data)
}

// ResponseJ (composing) returns a HTTP response with status and a json object as data
func (w WebPart) ResponseJ(status int, v interface{}) WebPart {
	return Compose(w, ResponseJ(status, v))
}

// OK returns 200 response with data
func OK(data []byte) WebPart {
	return Response(http.StatusOK, data)
}

// OK (composing) returns 200 response with data
func (w WebPart) OK(data []byte) WebPart {
	return Compose(w, OK(data))
}

// OKR returns 200 response with a reader containing the body
func OKR(reader func() io.ReadCloser) WebPart {
	return ResponseR(http.StatusOK, reader)
}

// OKR (composing) returns 200 response with a reader containing the body
func (w WebPart) OKR(reader func() io.ReadCloser) WebPart {
	return Compose(w, OKR(reader))
}

// OKS returns 200 response with text
func OKS(text string) WebPart {
	return ContentType(ContentTypeText).OK([]byte(text))
}

// OKS returns 200 response with text
func (w WebPart) OKS(text string) WebPart {
	return Compose(w, OKS(text))
}

// OKJ returns 200 response with data encoded as json
func OKJ(v interface{}) WebPart {
	return ContentType(ContentTypeJSON).ResponseJ(http.StatusOK, v)
}

// OKJ (composing) returns 200 response with data encoded as json
func (w WebPart) OKJ(v interface{}) WebPart {
	return Compose(w, OKJ(v))
}

//OKExtras returns 200 response with extras as json
func OKExtras() WebPart {
	return func(u WebUnit) *WebUnit {
		return OKJ(u.Extras())(u)
	}
}

//OKExtras returns 200 response with extras as json
func (w WebPart) OKExtras() WebPart {
	return Compose(w, OKExtras())
}

// Bad returns 400 response with error message
func Bad(data []byte) WebPart {
	return Response(http.StatusBadRequest, data)
}

// Bad (composing) returns 400 response with error message
func (w WebPart) Bad(data []byte) WebPart {
	return Compose(w, Bad(data))
}

// BadJ returns 400 response with json object as data
func BadJ(v interface{}) WebPart {
	return ContentType(ContentTypeJSON).ResponseJ(http.StatusBadRequest, v)
}

// BadJ (composing) returns 400 response with json object as data
func (w WebPart) BadJ(v interface{}) WebPart {
	return Compose(w, BadJ(v))
}

// NotFound [404]
func NotFound(data []byte) WebPart {
	return Response(http.StatusNotFound, data)
}

// NotFound [404]
func (w WebPart) NotFound(data []byte) WebPart {
	return Compose(w, NotFound(data))
}

// NotFoundS [404] with message
func NotFoundS(message string) WebPart {
	return ContentType(ContentTypeText).NotFound([]byte(message))
}

// NotFoundS (composing) [404] with message
func (w WebPart) NotFoundS(message string) WebPart {
	return Compose(w, NotFoundS(message))
}

// NotFoundJ [404] with json
func NotFoundJ(v interface{}) WebPart {
	return ContentType(ContentTypeJSON).ResponseJ(http.StatusNotFound, v)
}

// NotFoundJ (composing) [404] with json
func (w WebPart) NotFoundJ(v interface{}) WebPart {
	return Compose(w, NotFoundJ(v))
}

// RespondWith is an id function (non-modifing)
// Just for convenience
func (w WebPart) RespondWith() WebPart {
	return func(u WebUnit) *WebUnit { return &u }
}

// Error struct
type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	data, err := json.Marshal(e)
	Try(err)
	return string(data)
}

// Panic with HTTP status code and an error message
func Panic(status int, message string) WebPart {
	return func(u WebUnit) *WebUnit {
		panic(Error{Status: status, Message: message})
	}
}

// Panic (composing) with HTTP status code and an error message
func (w WebPart) Panic(status int, message string) WebPart {
	return Compose(w, func(u WebUnit) *WebUnit {
		panic(Error{Status: status, Message: message})
	})
}

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
