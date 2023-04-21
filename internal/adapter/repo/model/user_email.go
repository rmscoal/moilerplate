package model

type UserEmail struct {
	BaseModelId

	UserId    string
	Email     string `gorm:"unique"`
	IsPrimary bool

	BaseModelStamps
	BaseModelSoftDelete
}
