package model

type UserCredential struct {
	BaseModelId

	UserId   string
	Username string `gorm:"type:varchar(25);index:,unique,type:btree"`
	Password string `gorm:"type:varchar(255);index:,type:btree"`

	BaseModelStamps
	BaseModelSoftDelete
}
