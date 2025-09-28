package controllers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	db "github.com/SonikSeven/password-manager/db/sqlc"
	"github.com/SonikSeven/password-manager/schemas"
	"github.com/gin-gonic/gin"

	"github.com/SonikSeven/password-manager/util"
)

func GetUserID(c *gin.Context) (int64, error) {
	val, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("user ID not found")
	}
	id, ok := val.(int64)
	if !ok {
		return 0, errors.New("invalid user ID type")
	}
	return id, nil
}

type ListPasswordsResponse struct {
	Passwords []db.ListPasswordsRow `json:"passwords"`
}

type GetPasswordResponse struct {
	Password db.Password `json:"password"`
}

type PasswordController struct {
	config util.Config
	db     *db.Queries
	ctx    context.Context
}

func NewPasswordController(config util.Config, db *db.Queries, ctx context.Context) *PasswordController {
	return &PasswordController{config, db, ctx}
}

func (pc *PasswordController) ListPasswords(ctx *gin.Context) {
	userID, err := GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	search := ctx.Query("search")
	domain := ctx.Query("domain")

	passwords, err := pc.db.ListPasswords(ctx, db.ListPasswordsParams{
		UserID: userID,
		Search: search,
		Filter: domain,
	})

	if err != nil {
		log.Println("ListPasswords error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve passwords"})
		return
	}

	resp := make([]schemas.PasswordResponse, len(passwords))
	for i, p := range passwords {
		resp[i] = schemas.MapPasswordRow(p)
	}

	ctx.JSON(http.StatusOK, gin.H{"passwords": resp})
}

func (pc *PasswordController) GetPassword(ctx *gin.Context) {
	userID, err := GetUserID(ctx)
	passwordIDStr := ctx.Param("id")

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	passwordID, err := strconv.ParseInt(passwordIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	args := &db.GetPasswordByIDParams{
		ID:     int64(passwordID),
		UserID: userID,
	}

	password, err := pc.db.GetPasswordByID(ctx, *args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Password not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve password"})
		}
		log.Println("GetPassword error:", err)
		return
	}

	decryptedPassword, err := util.Decrypt(password.Password, pc.config.EncryptionKeyRaw)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decrypt password"})
		return
	}
	password.Password = decryptedPassword

	resp := schemas.MapPassword(password)

	ctx.JSON(http.StatusOK, gin.H{"password": resp})
}

func (pc *PasswordController) CreatePassword(ctx *gin.Context) {
	userID, err := GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var payload *schemas.CreatePassword

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
		return
	}

	encryptedPassword, err := util.Encrypt(payload.Password, pc.config.EncryptionKeyRaw)
	if err != nil {
		log.Println("CreatePassword error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt password"})
		return
	}

	args := &db.CreatePasswordParams{
		UserID:   userID,
		Service:  payload.Service,
		Username: payload.Username,
		Password: encryptedPassword,
		Url:      sql.NullString{String: payload.URL, Valid: payload.URL != ""},
		Notes:    sql.NullString{String: payload.Notes, Valid: payload.Notes != ""},
		Icon:     sql.NullString{String: payload.Icon, Valid: payload.Icon != ""},
	}

	_, err = pc.db.CreatePassword(ctx, *args)
	if err != nil {
		log.Println("CreatePassword error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create password"})
		return
	}

	ctx.Status(http.StatusCreated)
}

func (pc *PasswordController) UpdatePassword(ctx *gin.Context) {
	userID, err := GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	passwordIDStr := ctx.Param("id")
	var payload *schemas.UpdatePassword

	passwordID, err := strconv.ParseInt(passwordIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
		return
	}

	encryptedPassword, err := util.Encrypt(payload.Password, pc.config.EncryptionKeyRaw)
	if err != nil {
		log.Println("CreatePassword error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt password"})
		return
	}

	args := &db.UpdatePasswordParams{
		ID:       passwordID,
		UserID:   userID,
		Service:  payload.Service,
		Username: payload.Username,
		Password: encryptedPassword,
		Url:      sql.NullString{String: payload.URL, Valid: payload.URL != ""},
		Notes:    sql.NullString{String: payload.Notes, Valid: payload.Notes != ""},
		Icon:     sql.NullString{String: payload.Icon, Valid: payload.Icon != ""},
	}

	_, err = pc.db.UpdatePassword(ctx, *args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Password not found"})
			return
		}

		log.Println("UpdatePassword error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	ctx.Status(http.StatusOK)
}

func (pc *PasswordController) DeletePassword(ctx *gin.Context) {
	userID, err := GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	passwordIDStr := ctx.Param("id")

	passwordID, err := strconv.ParseInt(passwordIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	args := &db.DeletePasswordParams{
		ID:     passwordID,
		UserID: userID,
	}

	_, err = pc.db.DeletePassword(ctx, *args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Password not found"})
			return
		}

		log.Println("DeletePassword error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete password"})
		return
	}

	ctx.Status(http.StatusOK)
}
