package domain

import (
	"context"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
)

type User struct {
	Id          string
	FirstName   string
	LastName    string
	Emails      []vo.UserEmail
	PhoneNumber string
	Credential  vo.UserCredential
}

func (v User) ValidateWithContext(ctx context.Context) error {
	if err := validation.ValidateStructWithContext(ctx, &v,
		validation.Field(&v.Id, validation.When(v.Id != "", validation.Required, is.UUIDv4)),
		validation.Field(&v.FirstName, validation.Required, validation.Length(3, 20)),
		validation.Field(&v.LastName, validation.Required, validation.Length(3, 25)),
		validation.Field(&v.Credential),
		validation.Field(&v.Emails, validation.Each(validation.NotNil)),
		// see: https://stackoverflow.com/questions/44670612/regex-for-indonesian-phone-number
		validation.Field(&v.PhoneNumber, validation.Required, validation.Match(regexp.MustCompile(
			`(\+62 ((\d{3}([ -]\d{3,})([- ]\d{4,})?)|(\d+)))|(\(\d+\) \d+)|\d{3}( \d+)+|(\d+[ -]\d+)|\d+`),
		)),
	); err != nil {
		return err
	}

	for _, email := range v.Emails {
		err := email.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *User) AddEmail(newEmail vo.UserEmail) {
	for _, email := range u.Emails {
		if email.Equals(newEmail) {
			return
		}
	}
	u.Emails = append(u.Emails, newEmail)
}
