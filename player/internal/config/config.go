package config

import (
	"embed"
	"io/fs"
)

const CONFIG_BOT_THROTTLE = 100 // ms of sleep between each game
const CONFIG_DICTIONARY_FILENAME = "corncob_lowercase.txt"
const CONFIG_DICTIONARY_FILEPATH = "data/" + CONFIG_DICTIONARY_FILENAME
const CONFIG_GAME_WORDLENGTH = 5
const CONFIG_SERVER_URL = "http://localhost:8080"

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
