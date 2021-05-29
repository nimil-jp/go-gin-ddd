package validate

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/ja"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	jaTranslations "github.com/go-playground/validator/v10/translations/ja"
)

var (
	v, ok = binding.Validator.Engine().(*validator.Validate)
	uni   *ut.UniversalTranslator
	trans ut.Translator
)

func init() {
	jp := ja.New()
	uni = ut.New(jp, jp)
	trans, _ = uni.GetTranslator("ja")

	_ = jaTranslations.RegisterDefaultTranslations(v, trans)

	if ok {
		v.RegisterTagNameFunc(
			func(fld reflect.StructField) string {
				if value, ok := values[toSnakeCase(fld.Name)]; ok {
					return value
				}
				return toSnakeCase(fld.Name)
			},
		)

		registerAll()
	}
}

type errs struct {
	errors map[string][]*err
}

type err struct {
	Message string `json:"message"`
}

func NewErrs() *errs {
	te := new(errs)
	te.errors = make(map[string][]*err)
	return te
}

func (te *errs) Add(key string, value string) {
	if _, ok := te.errors[key]; !ok {
		te.errors[key] = []*err{{Message: value}}
	} else {
		te.errors[key] = append(te.errors[key], &err{Message: value})
	}
}

func (te *errs) Get() map[string][]*err {
	return te.errors
}

func BindAndValidate(c *gin.Context, request interface{}) bool {
	errs := NewErrs()
	err := c.BindJSON(request)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			c.JSON(http.StatusInternalServerError, err.Error())
			return false
		}
		if _, ok := err.(validator.ValidationErrors); !ok {
			c.JSON(http.StatusBadRequest, err.Error())
			return false
		}
		for _, f := range err.(validator.ValidationErrors) {
			errs.Add(toSnakeCase(f.StructField()), f.Translate(trans))
		}
	}
	if len(errs.Get()) != 0 {
		encoder := json.NewEncoder(c.Writer)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(map[string]interface{}{"errors": errs.Get()}); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return false
		}
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Status(http.StatusBadRequest)
		return false
	}
	return true
}

func register(tag string, fn validator.Func, translation string) {
	_ = v.RegisterValidation(tag, fn)
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
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}
	_ = v.RegisterTranslation(tag, trans, registrationFunc(tag, translation, true), translateFunc)
}

func toSnakeCase(str string) string {
	snake := regexp.MustCompile("(.)([A-Z][a-z]+)").ReplaceAllString(str, "${1}_${2}")
	snake = regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
