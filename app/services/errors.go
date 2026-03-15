package services

import "errors"

var (
	ErrPaymentDeclined    = errors.New("payment declined")
	ErrOrderNotFound      = errors.New("order not found")
	ErrMerchantNotFound   = errors.New("merchant not found")
	ErrInstallationExists = errors.New("installation already exists")
	ErrTokenFetch         = errors.New("failed to fetch token")
	ErrNilDependency      = errors.New("dependency is nil")
	ErrInvalidConfig      = errors.New("invalid configuration")
)
