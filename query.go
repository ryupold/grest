package grest

// Query checks Request for required keys.
// Pass through if requirements are met.
// Additionally all query parameters are put into Extras
func Query(requiredKeys ...string) WebPart {
	return QueryOrFail(nil, requiredKeys...)
}

// Query checks Request for required keys.
// Pass through if requirements are met.
// Additionally all query parameters are put into Extras
func (w WebPart) Query(requiredKeys ...string) WebPart {
	return Compose(w, Query(requiredKeys...))
}

// QueryOrFail checks Request for required keys.
// Pass through if requirements are met.
// Additionally all query parameters are put into Extras (only their first value though)
// Return errResult if not.
// It's ok to pass nil if this WebPart should be skipped on error
func QueryOrFail(errResult *WebPart, requiredKeys ...string) WebPart {
	return func(u WebUnit) *WebUnit {
		result := &u

		query := u.Request.URL.Query()
		for _, rk := range requiredKeys {
			if query.Get(rk) == "" {
				result = nil
				break
			}
		}

		if result == nil && errResult != nil {
			return (*errResult)(u)
		}

		if result != nil {
			for k, v := range query {
				if len(v) > 0 {
					result.PutExtra(k, v[0])
				} else {
					result.PutExtra(k, "")
				}
			}
		}

		return result
	}
}

// QueryOrFail checks Request for required keys.
// Pass through if requirements are met.
// Additionally all query parameters are put into Extras
// Return errResult if not.
// It's ok to pass nil if this WebPart should be skipped on error
func (w WebPart) QueryOrFail(errResult *WebPart, requiredKeys ...string) WebPart {
	return Compose(w, QueryOrFail(errResult, requiredKeys...))
}
