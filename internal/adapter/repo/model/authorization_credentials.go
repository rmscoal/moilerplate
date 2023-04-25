package model

import "time"

type AuthorizationCredential struct {
	BaseModelId

	Version int  `gorm:"default:1"`
	Issued  bool `gorm:"default:false"`

	ParentId *string `gorm:"default:null"`
	Parent   *AuthorizationCredential

	UserId   string
	IssuedAt time.Time `gorm:"default:now()"`
}
