package installation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pintarlabs/license-server/internal/database"
	"github.com/pintarlabs/license-server/internal/middleware"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) GetAll(c *gin.Context) {
	var installations []database.Installation
	if err := h.DB.Find(&installations).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch installations")
		return
	}
	middleware.SuccessResponse(c, installations)
}
