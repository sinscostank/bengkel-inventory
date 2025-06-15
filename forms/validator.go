package forms

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

//ValidateFullName implements validator.Func
func ValidateFullName(fl validator.FieldLevel) bool {
	//Remove the extra space
	space := regexp.MustCompile(`\s+`)
	name := space.ReplaceAllString(fl.Field().String(), " ")

	//Remove trailing spaces
	name = strings.TrimSpace(name)

	//To support all possible languages
	matched, _ := regexp.Match(`^[^±!@£$%^&*_+§¡€#¢§¶•ªº«\\/<>?:;'"|=.,0123456789]{3,20}$`, []byte(name))
	return matched
}