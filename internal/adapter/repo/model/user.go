package model

type User struct {
	BaseModelId

	FirstName                string
	LastName                 string
	PhoneNumber              string                    `gorm:"index:,unique,type:btree;size:30"`
	UserCredential           UserCredential            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	AuthorizationCredentials []AuthorizationCredential `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserEmails               []UserEmail               `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	BaseModelStamps
	BaseModelSoftDelete
}
