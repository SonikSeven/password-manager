package routes

import (
	"github.com/SonikSeven/password-manager/controllers"
	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	userController controllers.UserController
}

func NewRouteContact(userController controllers.UserController) UserRoutes {
	return UserRoutes{userController}
}

func (cr *UserRoutes) UserRoute(rg *gin.RouterGroup) {
	router := rg.Group("users")
	router.POST("/", cr.userController.CreateUser)
	router.GET("/", cr.userController.GetUserByID)
}
