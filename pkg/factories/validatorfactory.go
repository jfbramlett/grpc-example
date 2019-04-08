package factories

import (
	"fmt"
	"github.com/jfbramlett/grpc-example/pkg/rundef"
)

type ValidatorFactory interface {
	GetValidator(name string) (rundef.Validator, error)
}

func GetValidatorFactory() ValidatorFactory {
	validatorFactory := &simpleValidatorFactory{}
	validatorFactory.registerValidators()
	return validatorFactory
}


type simpleValidatorFactory struct {
	validators 		map[string]rundef.Validator
}


func (v *simpleValidatorFactory) GetValidator(name string) (rundef.Validator, error) {
	validator, found := v.validators[name]
	if !found {
		return nil, fmt.Errorf("failed to find validator named %s", name)
	}

	return validator, nil
}


func (v *simpleValidatorFactory) registerValidators() {
	v.validators = make(map[string]rundef.Validator)
	v.validators["default"] = &rundef.DefaultValidator{}
}
