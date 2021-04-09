package model

type Port struct {
	BaseUUIDModel
	IP         string `gorm:"type:varchar(64);"`
	Port       string `gorm:"type:varchar(12);"`
	State      string `gorm:"type:varchar(12);"`
	Protocol   string `gorm:"type:varchar(12);"`
	Retry      int    `gorm:"type:int;"`
	ServerType string `gorm:"type:varchar(64);"`
	Version    string `gorm:"type:varchar(64);"`
	Banner     string `gorm:"type:varchar(1024);"`
	IsSoft     bool   `gorm:"type:bool;"`
}

func (Port) TableName() string {
	return TabNamePort
}
