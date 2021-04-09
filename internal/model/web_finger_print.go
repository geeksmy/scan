package model

type Web struct {
	BaseModel
	Url         string `gorm:"type:varchar(64);"`
	StateCode   int    `gorm:"type:int"`
	Server      string `gorm:"type:varchar(64);"`
	Title       string `gorm:"type:varchar(128);"`
	FingerPrint string `gorm:"type:varchar(64);"`
	Retry       int    `gorm:"type:int"`
}

func (Web) TableName() string {
	return TabWeb
}
