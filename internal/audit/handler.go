package audit

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

func (h *Handler) GetAllLogs(c *gin.Context) {
	var logs []database.ValidationLog
	if err := h.DB.Order("created_at desc").Limit(100).Find(&logs).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch logs")
		return
	}
	middleware.SuccessResponse(c, logs)
}
