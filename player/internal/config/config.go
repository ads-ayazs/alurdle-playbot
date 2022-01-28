package config

import (
	"embed"
	"io/fs"
)

const (
	CONFIG_AWS_REGION = "ca-central-1"

	CONFIG_BOT_THROTTLE                            = 50 // ms of sleep between each game
	CONFIG_DICTIONARY_FILENAME                     = "corncob_lowercase.txt"
	CONFIG_DICTIONARY_FILEPATH                     = "data/" + CONFIG_DICTIONARY_FILENAME
	CONFIG_DICTIONARY_GENERATE_REAL_WORD_INJECTION = 0.3
	CONFIG_GAME_WORDLENGTH                         = 5
	CONFIG_SERVER_URL                              = "http://localhost:8080"
)

//go:embed data/*
var embFS embed.FS

func LoadEmbedFile(fp string) (fs.File, error) {
	if len(fp) < 1 {
		return nil, ErrFilepath
	}

	f, err := embFS.Open(fp)
	if err != nil {
		return nil, err
	}

	return f, nil
}
