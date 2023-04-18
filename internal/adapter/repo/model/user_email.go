package model

type UserEmail struct {
	BaseModelId

	UserId    string `gorm:"index:,unique,composite:unique_user_email,priority:1,type:btree"`
	Email     string `gorm:"index:,unique,composite:unique_user_email,priority:2,type:btree"`
	IsPrimary bool

	BaseModelStamps
	BaseModelSoftDelete
}
