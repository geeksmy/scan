package model

type Port struct {
	BaseUUIDModel
	Protocol string `gorm:"type:varchar(12);"`
	TargetIP string `gorm:"type:varchar(32);"`

	// 用于返回
	Services []PortService
}

func (Port) TableName() string {
	return TabNamePort
}

type PortService struct {
	BaseUUIDModel
	PortID string `gorm:"primary_key;type:varchar(36);"`
	Port   string `gorm:"type:varchar(5);"`
	Type   string `gorm:"type:varchar(32);"`
	Banner string `gorm:"type:varchar(128);"`
}

func (PortService) TableName() string {
	return TabNamePortService
}
