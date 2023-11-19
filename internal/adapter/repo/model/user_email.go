package model

type UserEmail struct {
	BaseModelId

	UserId    string
	Email     string `gorm:"type:varchar(50);uniqueIndex"`
	IsPrimary bool

	BaseModelStamps
	BaseModelSoftDelete
}
