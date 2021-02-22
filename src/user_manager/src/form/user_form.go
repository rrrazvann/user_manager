package form

import (
	"fmt"
	"user_manager/src/entity"
	"user_manager/src/lib"
)

var (
	ErrEmailValidationFailed = fmt.Errorf("email validation failed")
)

type UserForm struct {
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Nickname   string     `json:"nickname"`
	Password   string     `json:"password"`
	Email      string     `json:"email"`
	Country    string     `json:"country"`
}

func ValidateUserForm(form *UserForm) (error, *entity.User) {
	if len(form.FirstName) <= 0 {
		return getValidationErrorForField("FirstName"), nil
	}

	if len(form.LastName) <= 0 {
		return getValidationErrorForField("LastName"), nil
	}

	if len(form.Nickname) <= 0 {
		return getValidationErrorForField("Nickname"), nil
	}

	if len(form.Password) <= 0 {
		return getValidationErrorForField("Password"), nil
	}

	if len(form.Email) <= 0 {
		return getValidationErrorForField("Email"), nil
	}
	if !lib.IsEmailValid(form.Email) {
		return ErrEmailValidationFailed, nil
	}

	if len(form.Country) <= 0 {
		return getValidationErrorForField("Country"), nil
	}

	user := entity.User{
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Nickname:  form.Nickname,
		Password:  form.Password,
		Email:     form.Email,
		Country:   form.Country,
	}

	return nil, &user
}

func getValidationErrorForField(field string) error {
	return fmt.Errorf("field %s should be set", field)
}