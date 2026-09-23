package license

import (
	"fmt"
	"strconv"
	"strings"
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

	fmt.Printf("[DEBUG] POST /api/v1/license/activate called. LicenseKey: %s, HWID: %s\n", req.LicenseKey, req.MachineFingerprint)

	var lic database.License
	if err := h.DB.Preload("Customer").First(&lic, "license_key = ?", req.LicenseKey).Error; err != nil {
		fmt.Printf("[DEBUG] Activation failed: License key not found\n")
		middleware.ErrorResponse(c, http.StatusUnauthorized, "INVALID_LICENSE", "License key not found")
		return
	}

	if lic.Customer.ID != "" && lic.Customer.Status != "ACTIVE" {
		middleware.ErrorResponse(c, http.StatusForbidden, "CUSTOMER_INACTIVE", "Customer account is not active")
		return
	}

	if lic.Status != "PENDING" && lic.Status != "ACTIVE" {
		middleware.ErrorResponse(c, http.StatusForbidden, "LICENSE_INACTIVE", "License is not active")
		return
	}

	// Fetch plan features
	var planFeatures []database.PlanFeature
	h.DB.Preload("Feature").Where("plan_id = ?", lic.PlanID).Find(&planFeatures)

	featuresMap := make(map[string]interface{})
	for _, pf := range planFeatures {
		codeKey := strings.ToLower(pf.Feature.Code)
		if pf.Feature.DataType == "NUMBER" {
			if val, err := strconv.ParseFloat(pf.Value, 64); err == nil {
				featuresMap[codeKey] = val
			} else {
				featuresMap[codeKey] = 0
			}
		} else if pf.Feature.DataType == "BOOLEAN" {
			featuresMap[codeKey] = (pf.Value == "1" || strings.ToLower(pf.Value) == "true")
		} else {
			featuresMap[codeKey] = pf.Value
		}
	}

	now := time.Now()
	var installation database.Installation
	err := h.DB.First(&installation, "license_id = ? AND machine_fingerprint = ?", lic.ID, req.MachineFingerprint).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Check Max Devices
			maxDevices := 1.0 // Default
			if val, ok := featuresMap["max_devices"].(float64); ok {
				maxDevices = val
			}

			var activeInstalls int64
			h.DB.Model(&database.Installation{}).Where("license_id = ?", lic.ID).Count(&activeInstalls)

			if float64(activeInstalls) >= maxDevices {
				middleware.ErrorResponse(c, http.StatusForbidden, "MAX_DEVICES_REACHED", "Batas maksimal perangkat untuk lisensi ini telah tercapai.")
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

		if lic.Status == "PENDING" {
			lic.Status = "ACTIVE"
			if lic.ActivatedAt == nil {
				lic.ActivatedAt = &now
			}
			h.DB.Save(&lic)
		}
	}



	claims := token.LicenseTokenClaims{
		LicenseID:      lic.ID,
		ProductID:      lic.ProductID,
		CustomerID:     lic.CustomerID,
		PlanID:         lic.PlanID,
		InstallationID:     installation.InstallationID,
		MachineFingerprint: installation.MachineFingerprint,
		Features:       featuresMap,
	}

	if lic.ExpiresAt != nil {
		claims.ExpiresAt = lic.ExpiresAt.Unix()
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

type RequestTrialRequest struct {
	MachineFingerprint string `json:"machine_fingerprint" binding:"required"`
	ProductCode        string `json:"product_code" binding:"required"`
	AppVersion         string `json:"app_version"`
	Hostname           string `json:"hostname"`
	Platform           string `json:"platform"`
}

func (h *ClientHandler) RequestTrial(c *gin.Context) {
	var req RequestTrialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	// Check if this machine already has a trial license
	var existingInstallation database.Installation
	err := h.DB.Preload("License").First(&existingInstallation, "machine_fingerprint = ?", req.MachineFingerprint).Error
	if err == nil && existingInstallation.License != nil && existingInstallation.License.Status == "ACTIVE" {
		middleware.ErrorResponse(c, http.StatusForbidden, "TRIAL_ALREADY_USED", "Perangkat ini sudah pernah menggunakan trial.")
		return
	}

	// Create a default trial customer if doesn't exist
	var trialCustomer database.Customer
	if err := h.DB.Where("customer_code = ?", "TRIAL-USERS").First(&trialCustomer).Error; err != nil {
		trialCustomer = database.Customer{
			ID:           uuid.NewString(),
			CustomerCode: "TRIAL-USERS",
			Name:         "Trial Users",
			Status:       "ACTIVE",
		}
		h.DB.Create(&trialCustomer)
	}

	// Find the specific product requested by the client
	var product database.Product
	if err := h.DB.First(&product, "product_code = ? AND status = 'ACTIVE'", req.ProductCode).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusNotFound, "PRODUCT_NOT_FOUND", "Product not found or inactive")
		return
	}
	var plan database.Plan
	if err := h.DB.First(&plan, "product_id = ? AND code = ? AND status = 'ACTIVE'", product.ID, "PRO").Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "NO_PLAN_FOUND", "No active PRO plan available for this product")
		return
	}

	now := time.Now()
	expiresAt := now.Add(14 * 24 * time.Hour)
	
	lic := database.License{
		ID:          uuid.NewString(),
		LicenseKey:  "TRIAL-" + strings.ToUpper(uuid.NewString()[:8]),
		CustomerID:  trialCustomer.ID,
		ProductID:   product.ID,
		PlanID:      plan.ID,
		Status:      "ACTIVE",
		ActivatedAt: &now,
		ExpiresAt:   &expiresAt,
	}
	h.DB.Create(&lic)

	installation := database.Installation{
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

	lic.InstallationID = &installation.ID
	h.DB.Save(&lic)

	// Fetch plan features
	var planFeatures []database.PlanFeature
	h.DB.Preload("Feature").Where("plan_id = ?", lic.PlanID).Find(&planFeatures)

	featuresMap := make(map[string]interface{})
	for _, pf := range planFeatures {
		codeKey := strings.ToLower(pf.Feature.Code)
		if pf.Feature.DataType == "NUMBER" {
			if val, err := strconv.ParseFloat(pf.Value, 64); err == nil {
				featuresMap[codeKey] = val
			} else {
				featuresMap[codeKey] = 0
			}
		} else if pf.Feature.DataType == "BOOLEAN" {
			featuresMap[codeKey] = (pf.Value == "1" || strings.ToLower(pf.Value) == "true")
		} else {
			featuresMap[codeKey] = pf.Value
		}
	}

	claims := token.LicenseTokenClaims{
		LicenseID:          lic.ID,
		ProductID:          lic.ProductID,
		CustomerID:         lic.CustomerID,
		PlanID:             lic.PlanID,
		InstallationID:     installation.InstallationID,
		MachineFingerprint: installation.MachineFingerprint,
		Features:           featuresMap,
		ExpiresAt:          expiresAt.Unix(),
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
		"license_key":     lic.LicenseKey,
		"expires_at":      expiresAt,
	})
}


