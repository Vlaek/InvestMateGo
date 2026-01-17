package entity

type PortfolioHierarchy struct {
	ParentID string `gorm:"primaryKey;autoIncrement:false;constraint:OnDelete:CASCADE"`
	ChildID  string `gorm:"primaryKey;autoIncrement:false;constraint:OnDelete:CASCADE"`
}
