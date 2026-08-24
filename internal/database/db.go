package database

import (
	"fmt"
	"log"

	"github.com/pintarlabs/license-server/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(
		&Product{}, &Plan{}, &Feature{}, &PlanFeature{},
		&Customer{}, &License{}, &Installation{}, &ClientDevice{}, &AppVersion{},
		&AuditLog{}, &ActivationLog{}, &ValidationLog{}, &TransferLog{}, &User{},
	)
	return db
}
