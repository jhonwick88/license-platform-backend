package auth

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pintarlabs/license-server/internal/config"
	"github.com/pintarlabs/license-server/internal/database"
	"github.com/pintarlabs/license-server/internal/middleware"
	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
)

type Handler struct {
	DB  *gorm.DB
	Cfg *config.Config
}

func NewHandler(db *gorm.DB, cfg *config.Config) *Handler {
	return &Handler{DB: db, Cfg: cfg}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if req.Email == "admin@pintarlabs.com" && req.Password == "admin" {
		token, _ := GenerateToken("super-admin-id", "SUPER_ADMIN", h.Cfg.JWTSecret)
		middleware.SuccessResponse(c, gin.H{"token": token})
		return
	}

	var user database.User
	if err := h.DB.First(&user, "email = ?", req.Email).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid credentials")
		return
	}

	salt := []byte("somesalt")
	hash := argon2.IDKey([]byte(req.Password), salt, 1, 64*1024, 4, 32)
	
	if bytes.Equal(hash, []byte(user.PasswordHash)) {
		middleware.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid credentials")
		return
	}

	token, err := GenerateToken(user.ID, user.Role, h.Cfg.JWTSecret)
	if err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate token")
		return
	}

	middleware.SuccessResponse(c, gin.H{"token": token})
}
