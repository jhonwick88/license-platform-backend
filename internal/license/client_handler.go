package license

import (
	"crypto/rsa"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pintarlabs/license-server/internal/database"
	"github.com/pintarlabs/license-server/internal/middleware"
	"github.com/pintarlabs/license-server/internal/token"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientHandler struct {
	DB         *gorm.DB
	PrivateKey *rsa.PrivateKey
}

func NewClientHandler(db *gorm.DB, pk *rsa.PrivateKey) *ClientHandler {
	return &ClientHandler{DB: db, PrivateKey: pk}
}

type ActivateRequest struct {
	LicenseKey         string `json:"license_key" binding:"required"`
	MachineFingerprint string `json:"machine_fingerprint" binding:"required"`
	AppVersion         string `json:"app_version"`
	Hostname           string `json:"hostname"`
	Platform           string `json:"platform"`
}

func (h *ClientHandler) Activate(c *gin.Context) {
	var req ActivateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	var lic database.License
	if err := h.DB.First(&lic, "license_key = ?", req.LicenseKey).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusUnauthorized, "INVALID_LICENSE", "License key not found")
		return
	}

	if lic.Status != "PENDING" && lic.Status != "ACTIVE" {
		middleware.ErrorResponse(c, http.StatusForbidden, "LICENSE_INACTIVE", "License is not active")
		return
	}

	now := time.Now()
	var installation database.Installation
	err := h.DB.First(&installation, "license_id = ? AND machine_fingerprint = ?", lic.ID, req.MachineFingerprint).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if lic.InstallationID != nil && *lic.InstallationID != "" {
				middleware.ErrorResponse(c, http.StatusForbidden, "LICENSE_BOUND", "License is already bound to another machine")
				return
			}
			installation = database.Installation{
				ID:                 uuid.NewString(),
				LicenseID:          lic.ID,
				InstallationID:     uuid.NewString(),
				MachineFingerprint: req.MachineFingerprint,
				Platform:           req.Platform,
				Hostname:           req.Hostname,
				AppVersion:         req.AppVersion,
				Status:             "ACTIVE",
				LastSeenAt:         &now,
				LastServerTime:     &now,
			}
			h.DB.Create(&installation)

			lic.Status = "ACTIVE"
			lic.InstallationID = &installation.ID
			lic.ActivatedAt = &now
			h.DB.Save(&lic)
		} else {
			middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Database error")
			return
		}
	} else {
		installation.LastSeenAt = &now
		h.DB.Save(&installation)
	}

	claims := token.LicenseTokenClaims{
		LicenseID:      lic.ID,
		ProductID:      lic.ProductID,
		CustomerID:     lic.CustomerID,
		PlanID:         lic.PlanID,
		InstallationID: installation.InstallationID,
		Features:       map[string]interface{}{},
	}

	signedToken, err := token.SignLicenseToken(claims, h.PrivateKey)
	if err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to sign license token")
		return
	}

	middleware.SuccessResponse(c, gin.H{
		"token":           signedToken,
		"installation_id": installation.InstallationID,
		"status":          lic.Status,
	})
}

func (h *ClientHandler) Validate(c *gin.Context) {
	middleware.SuccessResponse(c, gin.H{"message": "Valid"})
}
