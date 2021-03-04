package filedriller

import "github.com/pkg/errors"

func e(err error) {
	if err != nil {
		errors.Cause(err)
	}
}
