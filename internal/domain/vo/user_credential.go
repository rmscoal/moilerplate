// vo is the shorthand for value objects.
package vo

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UserCredential struct {
	Username string
	Password string
	Token    string
}

func (v UserCredential) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Username, validation.Required, validation.Length(3, 20)),
		validation.Field(&v.Password, validation.Required),
	)
}
