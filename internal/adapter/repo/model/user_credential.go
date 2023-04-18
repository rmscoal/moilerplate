package model

type UserCredential struct {
	BaseModelId

	UserId   string
	Username string `gorm:"index:,unique,type:btree"`
	Password string `gorm:"index"`

	BaseModelStamps
	BaseModelSoftDelete
}
