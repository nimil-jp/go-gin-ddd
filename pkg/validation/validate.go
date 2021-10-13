package validation

import (
	"reflect"

	"github.com/go-playground/locales/ja"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	jaTranslations "github.com/go-playground/validator/v10/translations/ja"

	"go-gin-ddd/pkg/util"
)

var (
	validate   = validator.New()
	translator ut.Translator
	uni        *ut.UniversalTranslator
)

func init() {
	validate.UseActualTagWhenTranslate()

	jp := ja.New()
	uni = ut.New(jp, jp)
	translator, _ = uni.GetTranslator("ja")

	_ = jaTranslations.RegisterDefaultTranslations(validate, translator)

	validate.RegisterTagNameFunc(
		func(fld reflect.StructField) string {
			if value, ok := values[util.SnakeCase(fld.Name)]; ok {
				return value
			}
			return util.SnakeCase(fld.Name)
		},
	)

	registerAll()
}

func Validate() *validator.Validate {
	return validate
}

func Translator() ut.Translator {
	return translator
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
		t, err := ut.T(fe.ActualTag(), fe.Field())
		if err != nil {
			return "入力された値が正しくありません。"
		}
		return t
	}
	_ = validate.RegisterTranslation(tag, translator, registrationFunc(tag, translation, true), translateFunc)
}
