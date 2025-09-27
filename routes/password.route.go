package routes

import (
	"github.com/SonikSeven/password-manager/controllers"
	"github.com/SonikSeven/password-manager/middleware"
	"github.com/gin-gonic/gin"
)

type PasswordRoutes struct {
	passwordController controllers.PasswordController
}

func NewRoutePassword(passwordController controllers.PasswordController) PasswordRoutes {
	return PasswordRoutes{passwordController}
}

func (pr *PasswordRoutes) PasswordRoute(rg *gin.RouterGroup, jwtSecret []byte) {
	router := rg.Group("passwords")
	router.Use(middleware.AuthMiddleware(jwtSecret))
	router.GET("/", pr.passwordController.ListPasswords)
	router.GET("/:id", pr.passwordController.GetPassword)
	router.POST("/", pr.passwordController.CreatePassword)
	router.PATCH("/:id", pr.passwordController.UpdatePassword)
	router.DELETE("/:id", pr.passwordController.DeletePassword)
}
