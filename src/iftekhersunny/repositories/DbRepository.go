package repositories

import (
	"iftekhersunny/configs"
	"iftekhersunny/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Database repository struct
type DbRepository struct {
	//
}

// Init method
func (dbRepository DbRepository) init() *gorm.DB {

	db, err := gorm.Open("mysql", configs.DB_USERNAME+":"+configs.DB_PASSWORD+"@/"+configs.DB_NAME+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}

	// migrate the schema
	db.AutoMigrate(&models.Url{})

	return db
}
