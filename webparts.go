package grest

import "net/http"
import "context"

// WebUnit wraps a http.ResponseWriter & a http.Request
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
