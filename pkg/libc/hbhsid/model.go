package hbhsid

import (
	"time"
)

func NewBaseHBHSIDModel() BaseHBHSIDModel {
	now := time.Now()
	return BaseHBHSIDModel{
		CreateAt: &now,
		UpdateAt: &now,
	}
}

func NewBaseHBHSIDModelWithUint32(i uint32) BaseHBHSIDModel {
	now := time.Now()
	return BaseHBHSIDModel{
		ID:       New(i),
		CreateAt: &now,
		UpdateAt: &now,
	}
}

type BaseHBHSIDModel struct {
	ID       ID         `sql:"primary_key;type:SERIAL;"`
	CreateAt *time.Time `sql:"type:timestamptz;default:current_timestamp"`
	UpdateAt *time.Time `sql:"type:timestamptz;default:current_timestamp"`
	DeleteAt *time.Time `sql:"type:timestamptz;default:null"`
}
