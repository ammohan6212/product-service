package models

type Category struct {
	ID       uint     `gorm:"primaryKey" json:"id"`
	Name     string   `gorm:"not null" json:"name"`
	ImageURL string   `gorm:"not null" json:"image_url"`
	Products []Product `json:"products,omitempty"`
}
