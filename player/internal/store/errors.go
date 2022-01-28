package store

import "errors"

var (
	ErrInvalidStoreManagerName    = errors.New("invalid store manager name")
	ErrFailedToCreateStoreManager = errors.New("failed to create store manager")
	ErrOneSmMarshalFailure        = errors.New("marshaling failed")

	ErrAwsConfig = errors.New("unable to load AWS SDK config")
)
