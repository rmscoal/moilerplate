package repo

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func AddError(prev, new error) error {
	var err error
	if prev == nil {
		err = new
	} else {
		err = fmt.Errorf("%v; %w", prev, new)
	}

	return err
}

func translateGORMError(err error) error {
	switch {
	case errors.Is(err, gorm.ErrDuplicatedKey):
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return fmt.Errorf("duplicated value of %s exists", pgErr.ConstraintName)
		}
		return gorm.ErrDuplicatedKey
	case errors.Is(err, gorm.ErrRecordNotFound):
		return fmt.Errorf("the record you are looking for is not found")
	}

	return err
}
