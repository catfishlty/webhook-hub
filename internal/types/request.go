package types

type UserCreateRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateRequest struct {
	Username string `json:"username" binding:"required"`
}

type UserResetPasswordRequest struct {
	Password    string `json:"password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
