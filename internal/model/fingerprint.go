package model

import (
	"regexp"
)

type Match struct {
	BaseUUIDModel
	IsSoft      bool   `gorm:"type:varchar(128);"`
	Service     string `gorm:"type:varchar(128);"`
	Pattern     string `gorm:"type:varchar(128);"`
	VersionInfo string `gorm:"type:varchar(128);"`

	PatternCompiled *regexp.Regexp
}

func (Match) TableName() string {
	return TabNameMatch
}

type Probe struct {
	BaseUUIDModel
	Name     string
	Data     string
	Protocol string

	Ports        string
	SSLPorts     string
	Rarity       int
	Fallback     string
	TotalWaitMS  int
	TCPWrappedMS int

	Matchs []*Match
}

func (Probe) TableName() string {
	return TabNameProbe
}
