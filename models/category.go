package models

type Category struct {
    ID       uint     `gorm:"primaryKey"`
    Name     string   `gorm:"not null"`
    ImageURL string   `gorm:"not null"`
    Products []Product
}
