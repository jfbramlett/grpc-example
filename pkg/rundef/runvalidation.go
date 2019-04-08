package rundef

import "reflect"

type Validator interface {
	Validate(results []reflect.Value) (bool, err)
}



type DefaultValidator struct {}

func (d *DefaultValidator) Validate(result []reflect.Value) (bool, error) {
	if len(result) == 2 {
		if result[1].IsNil() {
			return true, nil
		} else {
			return false, result[1].Interface().(error)
		}
	}
	return true, nil
}