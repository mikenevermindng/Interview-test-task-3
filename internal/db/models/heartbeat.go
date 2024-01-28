package models

import (
	"time"
)

type HeartBeat struct {
	ID           int `gorm:"primaryKey;autoIncrement:true;unique"`
	ResponseTime int
	Uri          string
	Ssl          bool
	Timeout      bool
	Status       string
	Error        *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
