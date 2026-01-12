package custom

import "errors"

var (
	ErrInvalidFilePath = errors.New("Invalid file path")
	ErrWalMissing      = errors.New("WAL is missing")
)
