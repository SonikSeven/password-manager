package controllers

import (
	"context"
	"net/http"
	"time"

	db "github.com/SonikSeven/password-manager/db/sqlc"
	"golang.org/x/crypto/bcrypt"

	"github.com/SonikSeven/password-manager/auth"
	"github.com/SonikSeven/password-manager/schemas"
	"github.com/SonikSeven/password-manager/util"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	config util.Config
	db     *db.Queries
	ctx    context.Context
}

func NewUserController(config util.Config, db *db.Queries, ctx context.Context) *UserController {
	return &UserController{config, db, ctx}
}

func (uc *UserController) Register(ctx *gin.Context) {
	var payload *schemas.CreateUser

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "Failed hashing password"})
		return
	}

	now := time.Now()
	args := &db.CreateUserParams{
		Username:     payload.Username,
		Email:        payload.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	user, err := uc.db.CreateUser(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed create user", "error": err.Error()})
		return
	}

	accessToken, refreshToken, err := auth.GenerateTokens([]byte(uc.config.JWTSecret), user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "Failed generating tokens", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":        "User successfully registered",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (uc *UserController) Login(ctx *gin.Context) {
	var payload *schemas.LoginUser

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
		return
	}

	user, err := uc.db.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "Failed",
			"error":  "Invalid email or password",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "Failed",
			"error":  "Invalid email or password",
		})
		return
	}

	accessToken, refreshToken, err := auth.GenerateTokens([]byte(uc.config.JWTSecret), user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "Failed generating tokens", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":        "User successfully authenticated",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
