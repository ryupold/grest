package grest

import "context"

//panicKey to save an error in the context
const panicKey contextKey = "panic"

//Panic puts an error into the panic context of the WebUnit, marking it effectivly as failed but resulting WebUnit which would result in a 500 error when creating the response
//if the WebUnit has already a panic attack, the given error is discarded; you cannot panic while panicing
func (u *WebUnit) Panic(err error) {
	if err != nil && u.GetPanic() == nil {
		u.Context = context.WithValue(u.Context, panicKey, err)
	}
}

//GetPanic returns the panic in the panic context of the WebUnit or nil if there is none
func (u WebUnit) GetPanic() error {
	err, _ := u.Context.Value(panicKey).(error)
	return err
}

//Panic puts an error into the panic context of the WebUnit, marking it effectivly as failed but resulting WebUnit which would result in a 500 error when creating the response
func Panic(err error) WebPart {
	return func(u WebUnit) *WebUnit {
		u.Panic(err)
		return &u
	}
}

//Panic puts an error into the panic context of the WebUnit, marking it effectivly as failed but resulting WebUnit which would result in a 500 error when creating the response
func (w WebPart) Panic(err error) WebPart {
	return Compose(w, Panic(err))
}

//=== Recover =====================================================================================

//Recover removes the active panic in a WebUnit and returns it if there is any
func (u *WebUnit) Recover() error {
	err := u.GetPanic()
	u.Context = context.WithValue(u.Context, panicKey, nil)
	return err
}

//Recover calls the panicAttack handler if a panic was present, effectively removing the panic error
//If the panicAttack handler returns an error it is put in place of the old panic error, leaving the WebUnit still in panic
func Recover(panicAttack func(error) error) WebPart {
	return func(u WebUnit) *WebUnit {
		u.Panic(panicAttack(u.Recover()))
		return &u
	}
}

//Recover calls the panicAttack handler if a panic was present, effectively removing the panic error
//If the panicAttack handler returns an error it is put in place of the old panic error, leaving the WebUnit still in panic
func (w WebPart) Recover(panicAttack func(error) error) WebPart {
	return Compose(w, Recover(panicAttack))
}

//BadDream transforms a panic attack into a different WebPart only if there was a panic
//The actual error is forgotten
func BadDream(wakeUp WebPart) WebPart {
	return func(u WebUnit) *WebUnit {
		if err := u.Recover(); err != nil {
			return wakeUp(u)
		}
		return &u
	}
}
