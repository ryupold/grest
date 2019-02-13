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

//=== Extras ======================================================================================
type contextKey string

const extraKey contextKey = "extras"

//Extras returns a Data object that can be set with PutExtra(s) during routing
func (u WebUnit) Extras() Data {
	data, _ := u.Context.Value(extraKey).(Data)
	return data
}

//PutExtra puts given key:value into Extras
func (u *WebUnit) PutExtra(key string, value interface{}) {
	d := u.Extras()
	if d == nil {
		d = Data{}
	}
	d[key] = value
	u.Context = context.WithValue(u.Context, extraKey, d)
}

//PutExtras puts given Data into extras, merging with previous key:values
func (u *WebUnit) PutExtras(extras Data) {
	e := u.Extras()
	if e == nil {
		u.Context = context.WithValue(u.Context, extraKey, extras)
	} else {
		u.Context = context.WithValue(u.Context, extraKey, e.Union(extras))
	}
}

//SetExtras overrides the previous extras with given Data object
func (u *WebUnit) SetExtras(extras Data) {
	u.Context = context.WithValue(u.Context, extraKey, extras)
}
