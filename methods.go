package grest

import "net/http"

// GET filters for requests with this method
func GET() WebPart {
	return Method(http.MethodGet)
}

// POST filters for requests with this method
func POST() WebPart { return Method(http.MethodPost) }

// PUT filters for requests with this method
func PUT() WebPart { return Method(http.MethodPut) }

// DELETE filters for requests with this method
func DELETE() WebPart { return Method(http.MethodDelete) }

// OPTIONS filters for requests with this method
func OPTIONS() WebPart { return Method(http.MethodOptions) }

// PATCH filters for requests with this method
func PATCH() WebPart { return Method(http.MethodPatch) }

// Method filters for requests with the given method
var Method = func(method string) WebPart {
	return func(u WebUnit) *WebUnit {
		if u.Request.Method == method {
			return &u
		}
		return nil
	}
}
