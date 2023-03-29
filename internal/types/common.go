package types

type CommonError struct {
	Code int
	Msg  string
	Err  error
}

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
