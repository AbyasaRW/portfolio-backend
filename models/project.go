package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	Title       string `json:"title"`
	Stack       string `json:"stack"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Image       string `json:"image"`
}
