package feature

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pintarlabs/license-server/internal/database"
	"github.com/pintarlabs/license-server/internal/middleware"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) GetAll(c *gin.Context) {
	var features []database.Feature
	if err := h.DB.Find(&features).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch features")
		return
	}
	middleware.SuccessResponse(c, features)
}

func (h *Handler) Create(c *gin.Context) {
	var req database.Feature
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	req.ID = uuid.NewString()
	if req.Status == "" {
		req.Status = "ACTIVE"
	}

	if err := h.DB.Create(&req).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create feature")
		return
	}
	middleware.SuccessResponse(c, req)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req database.Feature
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	var feature database.Feature
	if err := h.DB.First(&feature, "id = ?", id).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "Feature not found")
		return
	}

	feature.Name = req.Name
	feature.Description = req.Description
	feature.DataType = req.DataType

	if err := h.DB.Save(&feature).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update feature")
		return
	}
	middleware.SuccessResponse(c, feature)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	
	// Delete associated PlanFeatures first to avoid constraint errors
	if err := h.DB.Where("feature_id = ?", id).Delete(&database.PlanFeature{}).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to clear associated plan features")
		return
	}

	if err := h.DB.Where("id = ?", id).Delete(&database.Feature{}).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete feature")
		return
	}
	middleware.SuccessResponse(c, gin.H{"message": "Feature deleted successfully"})
}
