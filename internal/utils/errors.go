package utils

import "fmt"

func AddError(prev, new error) error {
	var err error
	if prev == nil {
		err = new
	} else {
		err = fmt.Errorf("%v; %w", prev, new)
	}

	return err
}
