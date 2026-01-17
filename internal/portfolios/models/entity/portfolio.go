package entity

import (
	"time"
)

type Portfolio struct {
	ID                        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserId                    string    `gorm:"uniqueIndex;not null;size:255"`
	Name                      string    `gorm:"size:255"`
	IsComposite               bool      `gorm:"not null;default:false"`
	ApplyTaxesOnPaidDividends bool      `gorm:"not null;default:false"`
	DividendTaxPercent        float32   `gorm:"default:0"`
	CreatedDate               time.Time `gorm:"autoUpdateTime;default:CURRENT_TIMESTAMP"`
	UpdatedDate               time.Time `gorm:"autoUpdateTime"`
	HasToken                  bool      `gorm:"default:false"`
	Token                     string    `gorm:"size:10;default:''"`
	Currency                  string    `gorm:"size:3;default:'RUB'"`
	Note                      string    `gorm:"size:255"`
	IsHidden                  bool      `gorm:"default:false"`
}
