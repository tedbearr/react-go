package dto

type RegisterDTO struct {
	Username string `json:"username" form:"username" binding:"required" validate:"min:3"`
	Password string `json:"password" form:"password" binding:"required" validate:"min:6"`
}
