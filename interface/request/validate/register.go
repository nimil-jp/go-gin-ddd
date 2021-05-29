package validate

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func registerAll() {
	register("phone", phoneValidator(), "{0}が正しい電話番号ではありません")
	registerTrans("eqfield", "フィールドが一致していません")
}

func phoneValidator() validator.Func {
	var phoneRegex = regexp.MustCompile(`(^$|^\d{10}$|^\d{11}$|^\d{3}-\d{4}-\d{4}$|^\d{2}-\d{4}-\d{4}$|^\d{3}-\d{3}-\d{4}$|^\d{4}-\d{2}-\d{4}$|^\d{5}-\d-\d{4}$)`)
	return func(fl validator.FieldLevel) bool {
		return phoneRegex.MatchString(fl.Field().String())
	}
}
