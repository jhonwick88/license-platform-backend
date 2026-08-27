package plan

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
	var plans []database.Plan
	if err := h.DB.Preload("Product").Find(&plans).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch plans")
		return
	}
	middleware.SuccessResponse(c, plans)
}

type FeatureInput struct {
	FeatureID string `json:"feature_id"`
	Value     string `json:"value"`
}

type CreatePlanRequest struct {
	ProductID   string         `json:"product_id"`
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	Features    []FeatureInput `json:"features"`
}

func (h *Handler) Create(c *gin.Context) {
	var req CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	
	planID := uuid.NewString()
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}

	plan := database.Plan{
		ID:          planID,
		ProductID:   req.ProductID,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Status:      status,
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&plan).Error; err != nil {
			return err
		}

		for _, f := range req.Features {
			planFeature := database.PlanFeature{
				ID:        uuid.NewString(),
				PlanID:    planID,
				FeatureID: f.FeatureID,
				Value:     f.Value,
			}
			if err := tx.Create(&planFeature).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create plan and features")
		return
	}

	middleware.SuccessResponse(c, plan)
}
