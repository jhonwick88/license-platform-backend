package customer

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
	var customers []database.Customer
	if err := h.DB.Preload("Licenses").Preload("Licenses.Product").Preload("Licenses.Plan").Find(&customers).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch customers")
		return
	}
	middleware.SuccessResponse(c, customers)
}

func (h *Handler) Create(c *gin.Context) {
	var req database.Customer
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	req.ID = uuid.NewString()
	if req.Status == "" {
		req.Status = "ACTIVE"
	}

	if err := h.DB.Create(&req).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create customer")
		return
	}
	middleware.SuccessResponse(c, req)
}


func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req database.Customer
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	var customer database.Customer
	if err := h.DB.First(&customer, "id = ?", id).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "Customer not found")
		return
	}

	customer.Name = req.Name
	customer.Email = req.Email
	customer.CompanyName = req.CompanyName
	customer.Phone = req.Phone

	if err := h.DB.Save(&customer).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update customer")
		return
	}
	middleware.SuccessResponse(c, customer)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	
	// Ensure we delete licenses tied to this customer (or handle it properly, but here we do Hard Delete)
	// Actually, just deleting the customer and cascading or manually deleting licenses:
	var licenses []database.License
	h.DB.Where("customer_id = ?", id).Find(&licenses)
	for _, l := range licenses {
		h.DB.Where("license_id = ?", l.ID).Delete(&database.Installation{})
		h.DB.Delete(&l)
	}

	if err := h.DB.Where("id = ?", id).Delete(&database.Customer{}).Error; err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete customer")
		return
	}
	middleware.SuccessResponse(c, gin.H{"message": "Customer deleted successfully"})
}

