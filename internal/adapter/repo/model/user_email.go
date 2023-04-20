package model

type UserEmail struct {
	BaseModelId

	UserId    string
	Email     string `gorm:"index:,unique,type:btree"`
	IsPrimary bool

	BaseModelStamps
	BaseModelSoftDelete
}
