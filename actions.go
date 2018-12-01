package grest

import (
	"fmt"
)

//Do executes an action. if it returns an error. Terminates WebPart chain (nil)
func Do(action func(WebUnit) error) WebPart {
	return func(u WebUnit) *WebUnit {
		if err := action(u); err != nil {
			fmt.Println(err)
			return nil
		}
		return &u
	}
}

//Do executes an action. if it returns an error. Terminates WebPart chain (nil)
func (w WebPart) Do(action func(WebUnit) error) WebPart {
	return Compose(w, Do(action))
}
