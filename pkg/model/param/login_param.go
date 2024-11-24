package param

type LoginParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResetPasswordParam struct {
	Password    string `json:"password" binding:"required"`
	OldPassword string `json:"oldPassword"`
}
