package domain

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
	BaseID

	Name        string `gorm:"type:varchar(150);not null"`
	Username    string `gorm:"type:varchar(20);not null;index:,unique"`
	Email       string `gorm:"type:varchar(50);not null;index:,unique"`
	PhoneNumber string `gorm:"type:varchar(20);not null;index:,unique"`
	Password    string `gorm:"type:varchar(255);not null"`

	BaseStamps
	BaseSoftDelete
}

func (v User) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.ID, validation.When(v.ID != "", validation.Required, is.UUIDv4.Error("invalid id format"))),
		validation.Field(&v.Name, validation.Required, validation.Length(1, 150).Error("name too long, maximum of 150")),
		validation.Field(&v.Username, validation.Required, validation.Length(1, 25).Error("username too long, maximum of 25")),
		validation.Field(&v.Email, validation.Required, is.Email.Error("invalid email format")),
		// see: https://stackoverflow.com/questions/44670612/regex-for-indonesian-phone-number
		validation.Field(&v.PhoneNumber, validation.Required,
			validation.Match(regexp.MustCompile(
				`(\+62 ((\d{3}([ -]\d{3,})([- ]\d{4,})?)|(\d+)))|(\(\d+\) \d+)|\d{3}( \d+)+|(\d+[ -]\d+)|\d+`),
			).Error("invalid phone number format"),
		),
	)
}
