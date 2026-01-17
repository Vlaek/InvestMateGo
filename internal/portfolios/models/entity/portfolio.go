package entity

import (
	"time"
)

type Portfolio struct {
	ID                        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserId                    string    `gorm:"not null;index;constraint:OnDelete:CASCADE"`
	Name                      string    `gorm:"size:255"`
	IsComposite               bool      `gorm:"not null;default:false"`
	ApplyTaxesOnPaidDividends bool      `gorm:"not null;default:false"`
	DividendTaxPercent        float32   `gorm:"default:0"`
	HasToken                  bool      `gorm:"default:false"`
	Token                     string    `gorm:"size:10;default:''"`
	Currency                  string    `gorm:"size:3;default:'RUB'"`
	Note                      string    `gorm:"size:255"`
	IsHidden                  bool      `gorm:"default:false"`
	CreatedAt                 time.Time `gorm:"autoCreateTime"`
	UpdatedAt                 time.Time `gorm:"autoUpdateTime"`
}
