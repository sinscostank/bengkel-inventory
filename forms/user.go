package forms

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

//UserForm ...
type UserForm struct{}

//LoginForm ...
type LoginForm struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=50"`
}

// RegisterForm ...
type RegisterForm struct {
	Name     string `form:"name" json:"name" binding:"required,min=3,max=20" validate:"fullName"`
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=50"`
}

//Name ...
func (f UserForm) Name(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Masukkan nama"
		}
		return errMsg[0]
	case "min", "max":
		return "Panjang nama harus 3 hingga 20 karakter"
	case "fullName":
		return "Nama tidak boleh mengandung angka atau simbol"
	default:
		return "Terjadi kesalahan, silakan coba lagi nanti"
	}
}

//Email ...
func (f UserForm) Email(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Masukkan email"
		}
		return errMsg[0]
	case "min", "max", "email":
		return "Masukkan email yang valid"
	default:
		return "Terjadi kesalahan, silakan coba lagi nanti"
	}
}

//Password ...
func (f UserForm) Password(tag string) (message string) {
	switch tag {
	case "required":
		return "Masukkan password"
	case "min", "max":
		return "Panjang password minima; 8 hingga 50 karakter"
	case "eqfield":
		return "Password tidak cocok"
	default:
		return "Terjadi kesalahan, silakan coba lagi nanti"
	}
}

//Signin ...
func (f UserForm) Login(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:

		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Terjadi kesalahan, silakan coba lagi nanti"
		}

		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Email" {
				return f.Email(err.Tag())
			}
			if err.Field() == "Password" {
				return f.Password(err.Tag())
			}
		}

	default:
		return "Request tidak valid"
	}

	return "Terjadi kesalahan, silakan coba lagi nanti"
}

//Register ...
func (f UserForm) Register(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:

		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Terjadi kesalahan, silakan coba lagi nanti"
		}

		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Name" {
				return f.Name(err.Tag())
			}

			if err.Field() == "Email" {
				return f.Email(err.Tag())
			}

			if err.Field() == "Password" {
				return f.Password(err.Tag())
			}

		}
	default:
		return "Request tidak valid"
	}

	return "Terjadi kesalahan, silakan coba lagi nanti"
}
