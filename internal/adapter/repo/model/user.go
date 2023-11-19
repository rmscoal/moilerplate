package model

type User struct {
	BaseModelId

	FirstName                string                    `gorm:"type:varchar(50)"`
	LastName                 string                    `gorm:"type:varchar(50)"`
	PhoneNumber              string                    `gorm:"type:varchar(25);index:,unique,type:btree;size:30"`
	UserCredential           UserCredential            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	AuthorizationCredentials []AuthorizationCredential `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserEmails               []UserEmail               `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	BaseModelStamps
	BaseModelSoftDelete
}
