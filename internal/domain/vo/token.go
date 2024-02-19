package vo

import "github.com/rmscoal/moilerplate/internal/domain"

type AccessVersioning struct {
	JTI      string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ParentID *string `gorm:"type:uuid;default:null"`
	UserID   string  `gorm:"type:uuid;not null"`
	Version  int     `gorm:"type:integer"`
	User     domain.User
	Parent   *AccessVersioning
}

type Token struct {
	AccessToken  string
	RefreshToken string
}
