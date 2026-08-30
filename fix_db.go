package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:Bismillah1!@tcp(localhost:3306)/license_platform?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	constraints := []struct{ table, name string }{
		{"installations", "fk_licenses_installation"},
		{"products", "fk_licenses_product"},
		{"plans", "fk_licenses_plan"},
		{"customers", "fk_licenses_customer"},
		{"products", "fk_features_product"},
		{"plans", "fk_plan_features_plan"},
		{"features", "fk_plan_features_feature"},
	}

	for _, c := range constraints {
		err := db.Exec(fmt.Sprintf("ALTER TABLE %s DROP FOREIGN KEY %s;", c.table, c.name)).Error
		if err != nil {
			fmt.Printf("Ignoring error for %s.%s: %v\n", c.table, c.name, err)
		} else {
			fmt.Printf("Dropped %s.%s\n", c.table, c.name)
		}
	}
}
