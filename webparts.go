package grest

import "net/http"
import "context"

//WebUnit wraps a http.ResponseWriter, a http.Request and a Context
//In Case of a request a single WebUnit is passed from WebPart to WebPart probably mutating in between until a WebPart writes the final response or results in nil (terminating the WebPart sequence). Imagine it like a tree structure.
type WebUnit struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Context context.Context
}

// WebPart is a functinal approach for routing requests.
// If the input context matches the requirements of the func
// a non-nil (optionally modified) pointer to a HttpContext is
// returned. Otherwise nil.
type WebPart func(WebUnit) *WebUnit

//contextKey is used to save data in the WebUnit context
type contextKey string
