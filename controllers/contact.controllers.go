package controllers

import (
	"context"
	"net/http"
	"time"

	db "github.com/SonikSeven/password-manager/db/sqlc"

	"github.com/SonikSeven/password-manager/schemas"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	db  *db.Queries
	ctx context.Context
}

func NewUserController(db *db.Queries, ctx context.Context) *UserController {
	return &UserController{db, ctx}
}

// Create user handler
func (cc *UserController) CreateUser(ctx *gin.Context) {
	var payload *schemas.CreateUser

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
	}

	now := time.Now()
	args := &db.CreateUserParams{
		username:      payload.Username,
		email:         payload.Email,
		passwrod_hash: payload.Password,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	_, err := cc.db.CreateUser(ctx, *args)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving contact", "error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "successfully created user"})
}
