package repositories

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNilORM = errors.New("orm dependency is nil")
)

func isNotFoundErr(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
