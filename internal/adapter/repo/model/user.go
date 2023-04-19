package model

type User struct {
	BaseModelId

	FirstName      string
	LastName       string
	PhoneNumber    string
	UserCredential UserCredential `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserEmails     []UserEmail    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	BaseModelStamps
	BaseModelSoftDelete
}
