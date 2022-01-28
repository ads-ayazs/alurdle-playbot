package playerbot

import "errors"

var (
	ErrNilDictionary  = errors.New("dictionary is nil")
	ErrInvalidBotName = errors.New("invalid bot name")
	ErrFailedAttempt  = errors.New("failed attempt")
)
