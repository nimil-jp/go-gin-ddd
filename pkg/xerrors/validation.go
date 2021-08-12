package xerrors

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"go-ddd/pkg/util"
	"go-ddd/pkg/validation"
)

type Validation struct {
	errors map[string][]string
}

func NewValidation() *Validation {
	return &Validation{errors: map[string][]string{}}
}

func (verr Validation) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Errors map[string][]string `json:"errors"`
		}{
			Errors: verr.errors,
		},
	)
}

func (verr *Validation) Add(fieldName string, value string) {
	key := util.SnakeCase(fieldName)
	if _, ok := verr.errors[key]; !ok {
		verr.errors[key] = []string{value}
	} else {
		verr.errors[key] = append(verr.errors[key], value)
	}
}

func (verr Validation) Error() string {
	b, _ := json.Marshal(verr)
	return string(b)
}

func (verr *Validation) Validate(request interface{}) (ok bool) {
	err := validation.Validate().Struct(request)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); !ok {
			for _, f := range err.(validator.ValidationErrors) {
				verr.Add(util.SnakeCase(f.StructField()), f.Translate(validation.Translator()))
			}
		}
	}
	return verr.IsInValid()
}

func (verr Validation) IsInValid() bool {
	return len(verr.errors) > 0
}
