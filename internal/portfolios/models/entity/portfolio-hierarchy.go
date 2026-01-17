package entity

type PortfolioHierarchy struct {
	ParentID uint64 `gorm:"primaryKey;autoIncrement:false;constraint:OnDelete:CASCADE"`
	ChildID  uint64 `gorm:"primaryKey;autoIncrement:false;constraint:OnDelete:CASCADE"`
}
