// vo is the shorthand for value objects.
package vo

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserEmail struct {
	Email     string
	IsPrimary bool
}

func (v *UserEmail) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Email, validation.Required, is.Email.Error(fmt.Sprintf("%s is not a valid email", v.Email))),
		validation.Field(&v.IsPrimary, validation.Required),
	)
}

func (u *UserEmail) Equals(email UserEmail) bool {
	return u.Email == email.Email
}
