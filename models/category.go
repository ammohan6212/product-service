package models

type Category struct {
	ID       uint     `gorm:"primaryKey" json:"id"`
	Name     string   `json:"name"`
	ImageURL string   `json:"image_url"`
	Products []Product `gorm:"foreignKey:CategoryID" json:"products"` // Optional: to preload all products in a category
}
