package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Feature struct {
	Code     string
	Name     string
	DataType string
}

func main() {
	dsn := "root:Bismillah1!@tcp(localhost:3306)/license_platform?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var features []Feature
	db.Table("features").Find(&features)
	for _, f := range features {
		fmt.Printf("Code: %s, Name: %s, DataType: %s\n", f.Code, f.Name, f.DataType)
	}
}
