package database

import (
	"time"
)

type User struct {
	ID           string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	Email        string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Role         string    `gorm:"type:varchar(50);not null;default:'ADMIN'" json:"role"`
	Status       string    `gorm:"type:varchar(50);not null;default:'ACTIVE'" json:"status"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Licenses     []License `gorm:"foreignKey:CustomerID" json:"licenses,omitempty"`
}

type Product struct {
	ID          string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	ProductCode string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"product_code"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"type:varchar(50);not null;default:'ACTIVE'" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	PlanFeatures []PlanFeature `gorm:"foreignKey:PlanID" json:"plan_features,omitempty"`
}

type Plan struct {
	ID          string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	ProductID   string    `gorm:"type:varchar(36);not null" json:"product_id"` 
	Product     Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Code        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"type:varchar(50);not null;default:'ACTIVE'" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	PlanFeatures []PlanFeature `gorm:"foreignKey:PlanID" json:"plan_features,omitempty"`
}

type Feature struct {
	ID          string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	ProductID   string    `gorm:"type:varchar(36);not null" json:"product_id"` 
	Product     Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Code        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	DataType    string    `gorm:"type:varchar(50);not null" json:"data_type"`
	Status      string    `gorm:"type:varchar(50);not null;default:'ACTIVE'" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	PlanFeatures []PlanFeature `gorm:"foreignKey:PlanID" json:"plan_features,omitempty"`
}

type PlanFeature struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	PlanID    string    `gorm:"type:varchar(36);not null" json:"plan_id"` 
	Plan      Plan      `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
	FeatureID string    `gorm:"type:varchar(36);not null" json:"feature_id"` 
	Feature   Feature   `gorm:"foreignKey:FeatureID" json:"feature,omitempty"`
	Value     string    `gorm:"type:varchar(255);not null" json:"value"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Customer struct {
	ID           string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	CustomerCode string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"customer_code"`
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	CompanyName  string    `gorm:"type:varchar(255)" json:"company_name"`
	Email        string    `gorm:"type:varchar(255)" json:"email"`
	Phone        string    `gorm:"type:varchar(100)" json:"phone"`
	Address      string    `gorm:"type:text" json:"address"`
	Notes        string    `gorm:"type:text" json:"notes"`
	Status       string    `gorm:"type:varchar(50);not null;default:'ACTIVE'" json:"status"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Licenses     []License `gorm:"foreignKey:CustomerID" json:"licenses,omitempty"`
}

type License struct {
	ID               string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	LicenseKey       string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"license_key"`
	CustomerID       string     `gorm:"type:varchar(36);not null" json:"customer_id"` 
	Customer         Customer   `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	ProductID        string     `gorm:"type:varchar(36);not null" json:"product_id"`
	Product          Product    `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	PlanID           string     `gorm:"type:varchar(36);not null" json:"plan_id"`
	Plan             Plan       `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
	InstallationID   *string    `gorm:"type:varchar(36)" json:"installation_id"` 
	Installation     *Installation `json:"installation,omitempty"`
	Status           string     `gorm:"type:varchar(50);not null;default:'PENDING'" json:"status"`
	ActivatedAt      *time.Time `json:"activated_at"`
	LastValidationAt *time.Time `json:"last_validation_at"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type Installation struct {
	ID                 string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	LicenseID          string     `gorm:"type:varchar(36);not null" json:"license_id"` 
	License            *License   `gorm:"foreignKey:LicenseID" json:"-"`
	InstallationID     string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"installation_id"`
	MachineFingerprint string     `gorm:"type:varchar(255);not null" json:"machine_fingerprint"`
	Platform           string     `gorm:"type:varchar(100)" json:"platform"`
	Hostname           string     `gorm:"type:varchar(255)" json:"hostname"`
	AppVersion         string     `gorm:"type:varchar(100)" json:"app_version"`
	Status             string     `gorm:"type:varchar(50);not null;default:'ACTIVE'" json:"status"`
	LastSeenAt         *time.Time `json:"last_seen_at"`
	LastServerTime     *time.Time `json:"last_server_time"`
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type ClientDevice struct {
	ID                 string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	InstallationID     string    `gorm:"type:varchar(36);not null" json:"installation_id"`
	DeviceFingerprint  string    `gorm:"type:varchar(255);not null" json:"device_fingerprint"`
	DeviceName         string    `gorm:"type:varchar(255)" json:"device_name"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type AppVersion struct {
	ID          string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	ProductID   string    `gorm:"type:varchar(36);not null" json:"product_id"` 
	Product     Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Version     string    `gorm:"type:varchar(50);not null" json:"version"`
	ReleaseDate time.Time `json:"release_date"`
	IsCritical  bool      `gorm:"default:false" json:"is_critical"`
	DownloadUrl string    `gorm:"type:text" json:"download_url"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type AuditLog struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Action    string    `gorm:"type:varchar(255);not null" json:"action"`
	Entity    string    `gorm:"type:varchar(100);not null" json:"entity"`
	EntityID  string    `gorm:"type:varchar(36);not null" json:"entity_id"`
	UserID    string    `gorm:"type:varchar(36);not null" json:"user_id"`
	Details   string    `gorm:"type:text" json:"details"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type ActivationLog struct {
	ID                 string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	LicenseID          string    `gorm:"type:varchar(36);not null" json:"license_id"`
	InstallationID     string    `gorm:"type:varchar(36);not null" json:"installation_id"`
	MachineFingerprint string    `gorm:"type:varchar(255);not null" json:"machine_fingerprint"`
	IPAddress          string    `gorm:"type:varchar(100)" json:"ip_address"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type ValidationLog struct {
	ID             string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	LicenseID      string    `gorm:"type:varchar(36);not null" json:"license_id"`
	InstallationID string    `gorm:"type:varchar(36);not null" json:"installation_id"`
	IPAddress      string    `gorm:"type:varchar(100)" json:"ip_address"`
	IsValid        bool      `gorm:"not null" json:"is_valid"`
	FailureReason  string    `gorm:"type:text" json:"failure_reason"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type TransferLog struct {
	ID            string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	LicenseID     string    `gorm:"type:varchar(36);not null" json:"license_id"`
	OldCustomerID string    `gorm:"type:varchar(36);not null" json:"old_customer_id"`
	NewCustomerID string    `gorm:"type:varchar(36);not null" json:"new_customer_id"`
	AuthorizedBy  string    `gorm:"type:varchar(36);not null" json:"authorized_by"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}








