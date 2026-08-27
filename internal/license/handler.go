package license

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pintarlabs/license-server/internal/database"
	"github.com/pintarlabs/license-server/internal/middleware"
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

func generateSecureRandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}

func (h *AdminHandler) Create(c *gin.Context) {
	var req database.License
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	var product database.Product
	if err := h.DB.First(&product, "id = ?", req.ProductID).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid Product ID")
		return
	}

	var plan database.Plan
	if err := h.DB.First(&plan, "id = ?", req.PlanID).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid Plan ID")
		return
	}

	req.ID = uuid.NewString()

	prodCode := strings.ToUpper(product.ProductCode)
	planCode := strings.ToUpper(plan.Code)

	for i := 0; i < 3; i++ {
		rand1 := generateSecureRandomString(4)
		rand2 := generateSecureRandomString(4)
		rand3 := generateSecureRandomString(4)

		req.LicenseKey = fmt.Sprintf("PL-%s-%s-%s-%s-%s", prodCode, planCode, rand1, rand2, rand3)

		var count int64
		h.DB.Model(&database.License{}).Where("license_key = ?", req.LicenseKey).Count(&count)
		if count == 0 {
			break
		}
	}

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

func (h *AdminHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Where("license_id = ?", id).Delete(&database.Installation{}).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to clear associated installations")
		return
	}
	if err := h.DB.Where("id = ?", id).Delete(&database.License{}).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete license")
		return
	}
	middleware.SuccessResponse(c, gin.H{"message": "License deleted"})
}
