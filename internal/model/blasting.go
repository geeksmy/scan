package model

type Blasting struct {
	BaseUUIDModel
	IP       string `gorm:"type:varchar(32);"`
	Port     string `gorm:"type:varchar(12);"`
	Username string `gorm:"type:varchar(64);"`
	Password string `gorm:"type:varchar(64);"`
	Server   string `gorm:"type:varchar(64);"`
	Retry    int    `gorm:"type:int;"`
}

func (Blasting) TableName() string {
	return TabBlasting
}
