package dto

import (
	"github.com/aria3ppp/password"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func passwordValidationRules() []validation.Rule {
	return []validation.Rule{
		validation.Length(8, 32),
		password.IsPassword().
			Numbers(2).
			LowerLetters(2).
			UpperLetters(2).
			SpecialChars(1),
	}
}

type UserCreateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var _ validation.Validatable = UserCreateRequest{}

func (r UserCreateRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Email, is.EmailFormat),
		validation.Field(
			&r.Password,
			passwordValidationRules()...,
		),
	)
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var _ validation.Validatable = UserLoginRequest{}

func (r UserLoginRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Email, validation.Required, is.EmailFormat),
		validation.Field(
			&r.Password,
			passwordValidationRules()...,
		),
	)
}

// ---------------------------------------------------

type UserLogoutRequest struct {
	Token string `json:"token"`
}

var _ validation.Validatable = UserLogoutRequest{}

func (r UserLogoutRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(
			&r.Token,
			validation.Required,
			is.UUID,
		),
	)
}

type UserRefreshTokenRequest struct {
	Token string `json:"token"`
}

var _ validation.Validatable = UserRefreshTokenRequest{}

func (r UserRefreshTokenRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(
			&r.Token,
			validation.Required,
			is.UUID,
		),
	)
}

// --------------------------------------------

type IDParam struct {
	ID int `param:"id"`
}

var _ validation.Validatable = IDParam{}

func (r IDParam) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.ID,
			validation.Required,
			validation.Min(1),
		),
	)
}
