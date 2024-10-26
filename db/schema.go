package db

import (
	"gorm.io/gorm"
)

type RateObject struct {
	gorm.Model
	ImageUrl   string   `gorm:"index" json:"image_url"`
	Liked      *bool    `json:"liked"`
	ProviderID int      `json:"provider_id"`
	Provider   Provider `json:"provider"`
}

type Provider struct {
	gorm.Model
	Url  string `json:"url"`
	Name string `json:"name"`
}
