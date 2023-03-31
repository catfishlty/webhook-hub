package types

import validation "github.com/go-ozzo/ozzo-validation/v4"

type UserCreateRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (r UserCreateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.Required),
		validation.Field(&r.Password, validation.Required),
	)
}

type UserUpdateRequest struct {
	Username string `json:"username" binding:"required"`
}

func (r UserUpdateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.Required),
	)
}

type UserResetPasswordRequest struct {
	Password    string `json:"password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (r UserResetPasswordRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Password, validation.Required),
		validation.Field(&r.NewPassword, validation.Required),
	)
}
