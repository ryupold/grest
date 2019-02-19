package grest

import (
	"context"
	"net/http"
)

const statusKey contextKey = "status"

//Status writes the given status header if not already happend
func Status(statusCode int) WebPart {
	return func(u WebUnit) *WebUnit {
		u.Context = context.WithValue(u.Context, statusKey, statusCode)
		return &u
	}
}

//Status writes the given status header if not already happend
func (w WebPart) Status(statusCode int) WebPart {
	return Compose(w, Status(statusCode))
}

//GetStatus returns the status code if Status(CODE) was already called before. Otherwise returns 0
func (u WebUnit) GetStatus() int {
	statusCode, _ := u.Context.Value(statusKey).(int)
	return statusCode
}

//OK is a convinience call that sets the OK (200 status)
func OK() WebPart {
	return Status(http.StatusOK)
}

//OK is a convinience call that sets the OK (200 status)
func (w WebPart) OK() WebPart {
	return Compose(w, OK())
}

//BadRequest is a convinience call that sets the BadRequest (400 status)
func BadRequest() WebPart {
	return Status(http.StatusBadRequest)
}

//BadRequest is a convinience call that sets the BadRequest (400 status)
func (w WebPart) BadRequest() WebPart {
	return Compose(w, BadRequest())
}

//NotFound is a convinience call that sets the PageNotFound (404 status)
func NotFound() WebPart {
	return Status(http.StatusNotFound)
}

//NotFound is a convinience call that sets the PageNotFound (404 status)
func (w WebPart) NotFound() WebPart {
	return Compose(w, NotFound())
}
