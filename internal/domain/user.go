package domain

import (
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/rmscoal/moilerplate/internal/domain/vo"
)

type User struct {
	Id          string
	FirstName   string
	LastName    string
	Emails      []vo.UserEmail
	PhoneNumber string
	Credential  vo.UserCredential
}

func (v User) Validate() error {
	if err := validation.ValidateStruct(&v,
		validation.Field(&v.Id, validation.When(v.Id != "", validation.Required, is.UUIDv4)),
		validation.Field(&v.FirstName, validation.Required, validation.Length(3, 20)),
		validation.Field(&v.LastName, validation.Required, validation.Length(3, 25)),
		validation.Field(&v.Credential),
		// see: https://stackoverflow.com/questions/44670612/regex-for-indonesian-phone-number
		validation.Field(&v.PhoneNumber, validation.Required, validation.Match(regexp.MustCompile(
			`(\+62 ((\d{3}([ -]\d{3,})([- ]\d{4,})?)|(\d+)))|(\(\d+\) \d+)|\d{3}( \d+)+|(\d+[ -]\d+)|\d+`),
		)),
	); err != nil {
		return err
	}

	return v.ValidateEmails()
}

func (v User) ValidateEmails() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Emails, validation.Required, validation.By(func(value any) error {
			emails, ok := value.([]vo.UserEmail)
			if !ok {
				return fmt.Errorf("unrecognizable user's emails")
			}

			hash := make(map[string]bool, 0)
			primaryCount := 0

			for _, email := range emails {
				if email.IsPrimary {
					primaryCount++
				}

				if primaryCount > 1 {
					return fmt.Errorf("there should only be one primary email")
				}

				if _, found := hash[email.Email]; found {
					return fmt.Errorf("should be unique")
				}
				hash[email.Email] = true
			}

			return nil
		})),
	)
}
