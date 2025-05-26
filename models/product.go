package models

type Product struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	Name       string   `gorm:"not null" json:"name"`
	Price      float64  `json:"price"`
	Stock      int      `json:"stock"`
	ImageURL   string   `gorm:"not null" json:"image_url"`
	CategoryID uint     `json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"category,omitempty"`
}