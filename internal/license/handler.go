package license

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pintarlabs/license-server/internal/database"
	"github.com/pintarlabs/license-server/internal/middleware"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminHandler struct {
	DB *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{DB: db}
}

func (h *AdminHandler) GetAll(c *gin.Context) {
	var licenses []database.License
	if err := h.DB.Preload("Product").Preload("Customer").Preload("Plan").Preload("Installation").Find(&licenses).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch licenses")
		return
	}
	middleware.SuccessResponse(c, licenses)
}

func (h *AdminHandler) Create(c *gin.Context) {
	var req database.License
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	req.ID = uuid.NewString()
	req.LicenseKey = "PL-" + uuid.NewString()[:8]
	if req.Status == "" {
		req.Status = "PENDING"
	}

	if err := h.DB.Create(&req).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate license")
		return
	}
	middleware.SuccessResponse(c, req)
}

func (h *AdminHandler) Suspend(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Model(&database.License{}).Where("id = ?", id).Update("status", "SUSPENDED").Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to suspend license")
		return
	}
	middleware.SuccessResponse(c, gin.H{"message": "License suspended"})
}

func (h *AdminHandler) Resume(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Model(&database.License{}).Where("id = ?", id).Update("status", "ACTIVE").Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to resume license")
		return
	}
	middleware.SuccessResponse(c, gin.H{"message": "License resumed"})
}

func (h *AdminHandler) Revoke(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Model(&database.License{}).Where("id = ?", id).Update("status", "REVOKED").Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to revoke license")
		return
	}
	middleware.SuccessResponse(c, gin.H{"message": "License revoked"})
}

