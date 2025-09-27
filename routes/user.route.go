package routes

import (
	"github.com/SonikSeven/password-manager/controllers"
	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	userController controllers.UserController
}

func NewRouteUser(userController controllers.UserController) UserRoutes {
	return UserRoutes{userController}
}

func (cr *UserRoutes) UserRoute(rg *gin.RouterGroup) {
	router := rg.Group("users")
	router.POST("/register", cr.userController.Register)
	router.POST("/login", cr.userController.Login)
}
