package models

import (
	"time"
)

type UserRequest struct {
	ID         int `gorm:"primaryKey;autoIncrement:true;unique"`
	Service    string
	Method     string
	Path       string
	Code       int
	ClientIp   string
	ReceivedAt time.Time
	DurationMs int
	CreatedAt  time.Time
}
