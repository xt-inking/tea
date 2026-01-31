package teavalid

import (
	"github.com/tea-frame-go/tea/teatypes"
)

type required interface {
	Validate(k teatypes.ToStringer) (err Error)
}

func Required() requiredRule {
	return requiredRule{}
}

type requiredRule struct{}

func (requiredRule) Validate(k teatypes.ToStringer) (err Error) {
	err = NewError(k.ToString(), "is required")
	return
}

func Optional() optionalRule {
	return optionalRule{}
}

type optionalRule struct{}

func (optionalRule) Validate(k teatypes.ToStringer) (err Error) {
	return
}
