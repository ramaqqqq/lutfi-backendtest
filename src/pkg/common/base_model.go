package common

import "time"

type BaseTime struct {
	Created time.Time `gorm:"default:current_timestamp" json:"created"`
	Updated time.Time `gorm:"default:current_timestamp" json:"updated"`
}
