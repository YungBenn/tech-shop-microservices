package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID        string    `gorm:"type:varchar(36);primaryKey;" json:"id"`
	Title     string    `gorm:"type:varchar(255)" json:"title"`
	Price     string    `gorm:"type:varchar(255)" json:"price"`
	Tag       []string  `gorm:"type:varchar(255)[]" json:"tag"`
	Discount  string    `gorm:"type:varchar(255)" json:"discount"`
	Image     []string  `gorm:"type:varchar(255)[]" json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.NewString()
	return
}
