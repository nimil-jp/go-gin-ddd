package validate

import (
	"encoding/json"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/locales/ja"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	jaTranslations "github.com/go-playground/validator/v10/translations/ja"
)

var (
	validate = validator.New()
	uni      *ut.UniversalTranslator
	trans    ut.Translator
)

func init() {
	jp := ja.New()
	uni = ut.New(jp, jp)
	trans, _ = uni.GetTranslator("ja")

	_ = jaTranslations.RegisterDefaultTranslations(validate, trans)

	validate.RegisterTagNameFunc(
		func(fld reflect.StructField) string {
			if value, ok := values[toSnakeCase(fld.Name)]; ok {
				return value
			}
			return toSnakeCase(fld.Name)
		},
	)

	registerAll()
}

type ValidationError struct {
	errors map[string][]string
}

func NewValidationError() *ValidationError {
	te := new(ValidationError)
	te.errors = make(map[string][]string)
	return te
}

func (verr ValidationError) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Errors map[string][]string `json:"errors"`
		}{
			Errors: verr.errors,
		},
	)
}

func (verr *ValidationError) Add(key string, value string) {
	if _, ok := verr.errors[key]; !ok {
		verr.errors[key] = []string{value}
	} else {
		verr.errors[key] = append(verr.errors[key], value)
	}
}

func (verr ValidationError) Error() string {
	b, _ := json.Marshal(verr)
	return string(b)
}

func (verr ValidationError) Validate(request interface{}) bool {
	err := validate.Struct(request)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); !ok {
			for _, f := range err.(validator.ValidationErrors) {
				verr.Add(toSnakeCase(f.StructField()), f.Translate(trans))
			}
		}
	}
	if len(verr.errors) > 0 {
		return false
	}
	return true
}

func register(tag string, fn validator.Func, translation string) {
	_ = validate.RegisterValidation(tag, fn)
	registerTrans(tag, translation)
}

func registerTrans(tag string, translation string) {
	registrationFunc := func(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
		return func(ut ut.Translator) (err error) {
			if err = ut.Add(tag, translation, override); err != nil {
				return
			}
			return
		}
	}

	translateFunc := func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field())
		if err != nil {
			log.Printf("warning: error translating FieldError: %#validate", fe)
			return fe.(error).Error()
		}
		return t
	}
	_ = validate.RegisterTranslation(tag, trans, registrationFunc(tag, translation, true), translateFunc)
}

func toSnakeCase(str string) string {
	snake := regexp.MustCompile("(.)([A-Z][a-z]+)").ReplaceAllString(str, "${1}_${2}")
	snake = regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
