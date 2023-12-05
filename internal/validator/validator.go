package validator

import (
	"slices"
	"strings"
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{
		Errors: map[string]string{},
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, ok := v.Errors[key]; !ok {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func NotEmpty(s string) bool {
	return strings.Trim(s, " ") != ""
}

func PremittedValues[K comparable](value K, permittedValues []K) bool {
	return slices.Contains(permittedValues, value)
}
