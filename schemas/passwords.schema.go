package schemas

import (
	"time"

	db "github.com/SonikSeven/password-manager/db/sqlc"
)

type CreatePassword struct {
	Service  string `json:"service" binding:"required,min=1"`
	Username string `json:"username" binding:"required,min=1"`
	Password string `json:"password" binding:"required,min=8"`
	URL      string `json:"url"`
	Notes    string `json:"notes"`
	Icon     string `json:"icon"`
}

type UpdatePassword struct {
	Service  string `json:"service" binding:"required,min=1"`
	Username string `json:"username" binding:"required,min=1"`
	Password string `json:"password" binding:"required,min=8"`
	URL      string `json:"url"`
	Notes    string `json:"notes"`
	Icon     string `json:"icon"`
}

type PasswordResponse struct {
	ID        int64   `json:"id"`
	UserID    int64   `json:"user_id"`
	Service   string  `json:"service"`
	Url       *string `json:"url,omitempty"`
	Notes     *string `json:"notes,omitempty"`
	Icon      *string `json:"icon,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

func MapPassword(p db.Password) PasswordResponse {
	var url, notes, icon *string

	if p.Url.Valid {
		url = &p.Url.String
	}
	if p.Notes.Valid {
		notes = &p.Notes.String
	}
	if p.Icon.Valid {
		icon = &p.Icon.String
	}

	return PasswordResponse{
		ID:        p.ID,
		UserID:    p.UserID,
		Service:   p.Service,
		Url:       url,
		Notes:     notes,
		Icon:      icon,
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
		UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
	}
}

func MapPasswordRow(p db.ListPasswordsRow) PasswordResponse {
	var url, notes, icon *string

	if p.Url.Valid {
		url = &p.Url.String
	}
	if p.Notes.Valid {
		notes = &p.Notes.String
	}
	if p.Icon.Valid {
		icon = &p.Icon.String
	}

	return PasswordResponse{
		ID:        p.ID,
		UserID:    p.UserID,
		Service:   p.Service,
		Url:       url,
		Notes:     notes,
		Icon:      icon,
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
		UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
	}
}
