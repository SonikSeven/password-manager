package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/SonikSeven/password-manager/controllers"
	dbCon "github.com/SonikSeven/password-manager/db/sqlc"
	"github.com/SonikSeven/password-manager/routes"
	"github.com/SonikSeven/password-manager/util"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	server *gin.Engine
	db     *dbCon.Queries
	ctx    context.Context

	UserController controllers.UserController
	UserRoutes     routes.UserRoutes
)

func init() {
	ctx = context.TODO()
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatalf("Could not loadconfig: %v", err)
	}

	conn, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	db = dbCon.New(conn)

	fmt.Println("PostgreSql connected successfully...")

	UserController = *controllers.NewUserController(db, ctx)
	UserRoutes = routes.NewRouteUser(UserController)

	server = gin.Default()
}

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	router := server.Group("/api")

	router.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "The password manager APi is working fine"})
	})

	UserRoutes.UserRoute(router)

	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": fmt.Sprintf("The specified route %s not found", ctx.Request.URL)})
	})

	log.Fatal(server.Run(":" + config.ServerAddress))
}
