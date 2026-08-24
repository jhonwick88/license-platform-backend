package main

import (
	"log"
	"github.com/gin-gonic/gin"
	
	"github.com/pintarlabs/license-server/internal/config"
	"github.com/pintarlabs/license-server/internal/database"
	"github.com/pintarlabs/license-server/internal/token"
	
	"github.com/pintarlabs/license-server/internal/auth"
	"github.com/pintarlabs/license-server/internal/customer"
	"github.com/pintarlabs/license-server/internal/feature"
	"github.com/pintarlabs/license-server/internal/license"
	"github.com/pintarlabs/license-server/internal/plan"
	"github.com/pintarlabs/license-server/internal/product"
	"github.com/pintarlabs/license-server/internal/installation"
	"github.com/pintarlabs/license-server/internal/audit"
)

func main() {
	cfg := config.Load()
	
	privateKey, err := token.GenerateRSAKeyPair()
	if err != nil {
		log.Fatalf("Failed to generate RSA keys: %v", err)
	}

	r := gin.Default()
	
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	var db = database.Connect(cfg)

	authHandler := auth.NewHandler(db, cfg)
	productHandler := product.NewHandler(db)
	planHandler := plan.NewHandler(db)
	featureHandler := feature.NewHandler(db)
	customerHandler := customer.NewHandler(db)
	licenseAdminHandler := license.NewAdminHandler(db)
	licenseClientHandler := license.NewClientHandler(db, privateKey)
	installationHandler := installation.NewHandler(db)
	auditHandler := audit.NewHandler(db)

	api := r.Group("/api/v1")
	{
		api.POST("/admin/login", authHandler.Login)
		
		api.POST("/license/activate", licenseClientHandler.Activate)
		api.POST("/license/validate", licenseClientHandler.Validate)

		admin := api.Group("/admin")
		admin.Use(auth.RequireAuth(cfg))
		{
			admin.GET("/products", productHandler.GetAll)
			admin.POST("/products", productHandler.Create)
			
			admin.GET("/plans", planHandler.GetAll)
			admin.POST("/plans", planHandler.Create)
			
			admin.GET("/features", featureHandler.GetAll)
			admin.POST("/features", featureHandler.Create)
			
			admin.GET("/customers", customerHandler.GetAll)
			admin.POST("/customers", customerHandler.Create)
			
			admin.GET("/licenses", licenseAdminHandler.GetAll)
			admin.POST("/licenses", licenseAdminHandler.Create)
			admin.POST("/licenses/:id/suspend", licenseAdminHandler.Suspend)
			admin.POST("/licenses/:id/resume", licenseAdminHandler.Resume)
			admin.POST("/licenses/:id/revoke", licenseAdminHandler.Revoke)

			admin.GET("/installations", installationHandler.GetAll)
			admin.GET("/logs", auditHandler.GetAllLogs)
		}
	}

	log.Printf("License Server running on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
