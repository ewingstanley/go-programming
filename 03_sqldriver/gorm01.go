package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type fabaodata struct{
	Title string 
}

func main() {

	dsn := "root:123456@tcp(localhost:3306)/dbtest"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println(db)
	// 迁移 schema
	// db.AutoMigrate(&Product{})

	// // Create
	var fineuserList []fabaodata
	db.Find(&fineuserList)
	fmt.Println(fineuserList)

}
