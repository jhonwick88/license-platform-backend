package product

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
	var products []database.Product
	if err := h.DB.Find(&products).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch products")
		return
	}
	middleware.SuccessResponse(c, products)
}

func (h *Handler) Create(c *gin.Context) {
	var req database.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	req.ID = uuid.NewString()
	if req.Status == "" {
		req.Status = "ACTIVE"
	}

	if err := h.DB.Create(&req).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create product")
		return
	}
	middleware.SuccessResponse(c, req)
}

