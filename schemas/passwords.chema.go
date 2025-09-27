package schemas

type CreatePassword struct {
	Name     string `json:"name" binding:"required,min=1"`
	Password string `json:"password" binding:"required,min=8"`
}

type UpdatePassword struct {
	Name     string `json:"name" binding:"required,min=1"`
	Password string `json:"password" binding:"required,min=8"`
}
