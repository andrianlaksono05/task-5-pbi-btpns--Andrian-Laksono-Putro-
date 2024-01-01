package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func KoneksiDatabase() {
	db, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/pbi_final"))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{})

	DB = db
}
