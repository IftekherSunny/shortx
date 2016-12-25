package models

import "github.com/jinzhu/gorm"

// Url model
type Url struct {
	gorm.Model
	LongUrl  string
	ShortUrl string
}
