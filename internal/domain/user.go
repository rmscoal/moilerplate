package domain

import (
	"context"
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/utils"
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
		validation.Field(&v.Emails, validation.Required, validation.Each(validation.NotNil)),
		// see: https://stackoverflow.com/questions/44670612/regex-for-indonesian-phone-number
		validation.Field(&v.PhoneNumber, validation.Required, validation.Match(regexp.MustCompile(
			`(\+62 ((\d{3}([ -]\d{3,})([- ]\d{4,})?)|(\d+)))|(\(\d+\) \d+)|\d{3}( \d+)+|(\d+[ -]\d+)|\d+`),
		)),
	); err != nil {
		return err
	}

	if err := v.ValidateEmailsWithContext(ctx); err != nil {
		return err
	}

	return nil
}

func (v User) ValidateEmailsWithContext(ctx context.Context) (err error) {
	for _, email := range v.Emails {
		vErr := email.Validate()
		if vErr != nil {
			err = utils.AddError(err, vErr)
		}
	}
	return err
}

func (u *User) VerifyOnePrimaryEmail() (err error) {
	counter := 0
	for _, email := range u.Emails {
		if email.IsPrimary {
			counter++
		}
	}
	if counter != 1 {
		err = fmt.Errorf("there exists multiple/none primary email")
	}
	return err
}

func (u *User) ReplaceEmails(emails []vo.UserEmail) {
	u.Emails = emails
}
