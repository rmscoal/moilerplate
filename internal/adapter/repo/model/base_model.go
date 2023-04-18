package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModelId struct {
	Id string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
}

type BaseModelStamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BaseModelSoftDelete struct {
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
