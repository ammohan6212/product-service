package models

type Product struct {
    ID         uint    `gorm:"primaryKey"`
    Name       string  `gorm:"not null"`
    Price      float64
    Stock      int
    CategoryID uint
}
