package validator

import (
	"fmt"
)

type ValidatorFactory interface {
	GetValidator(name string) (Validator, error)
}

func GetValidatorFactory() ValidatorFactory {
	validatorFactory := &simpleValidatorFactory{}
	validatorFactory.registerValidators()
	return validatorFactory
}


type simpleValidatorFactory struct {
	validators 		map[string]Validator
}


func (v *simpleValidatorFactory) GetValidator(name string) (Validator, error) {
	validator, found := v.validators[name]
	if !found {
		return nil, fmt.Errorf("failed to find validator named %s", name)
	}

	return validator, nil
}


func (v *simpleValidatorFactory) registerValidators() {
	v.validators = make(map[string]Validator)
	v.validators["default"] = &DefaultValidator{}
}
