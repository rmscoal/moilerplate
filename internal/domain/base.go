package domain

import (
	"time"

	"gorm.io/gorm"
)

type BaseID struct {
	ID string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
}

type BaseStamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BaseSoftDelete struct {
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
