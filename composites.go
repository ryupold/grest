package grest

// Choose takes a list of WebParts.
// On evaluation the list is tested from first to last
// and the first WebPart that statisfies is executed and returned
func Choose(options ...WebPart) WebPart {
	return func(unit WebUnit) *WebUnit {
		for _, c := range options {
			result := c(unit)
			if result != nil {
				return result
			}
		}
		return nil
	}
}

// Compose glues a list of WebParts togheter so that
// they are later evaluated from first to last in a row
// piping the result of each previous request into the next
// If one of them results in 'nil',
// nil is also the end result of the whole WebPart
func Compose(parts ...WebPart) WebPart {
	return func(unit WebUnit) *WebUnit {
		var result *WebUnit
		next := unit
		for _, p := range parts {
			if p == nil {
				break
			}
			result = p(next)
			if result == nil {
				break
			}

			next = *result
		}
		return result
	}
}
