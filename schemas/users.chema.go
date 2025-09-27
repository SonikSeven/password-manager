package schemas

type CreateUser struct {
	Username string `json:"username" binding:"required,min=1"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=1"`
}
