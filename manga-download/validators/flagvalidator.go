package validators

import (
	"errors"
	"strings"

	"github.com/xortock/mangafire-download/constants"
)

type TypeFlag struct {
	Value string
}

func ValidateTypeFlag(flag TypeFlag) error {
	if strings.ToLower(flag.Value) != constants.FILE_TYPE_ZIP && strings.ToLower(flag.Value) != constants.FILE_TYPE_CBZ {
		return errors.New(flag.Value + " file type not supported")
	}

	return nil
}

type DivisionFlag struct {
	Value string
}

func ValidateDivisionFlag(flag DivisionFlag) error {
	if strings.ToLower(flag.Value) != constants.DIVISION_CHAPTER && strings.ToLower(flag.Value) != constants.DVISION_VOLUME {
		return errors.New(flag.Value + " division type not supported")
	}
	return nil
}
