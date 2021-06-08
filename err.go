package filedriller

import (
	"github.com/pkg/errors"
	"os"
)

func e(err error) {
	if err != nil {
		errors.Cause(err)
		os.Exit(1)
	}
}
