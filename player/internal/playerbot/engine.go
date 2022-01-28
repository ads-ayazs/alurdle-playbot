package playerbot

import (
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	"aluance.io/wordleplayer/internal/config"
)

type GameEngine struct{}

func GetGameEngine() GameEngine {
	return GameEngine{}
}

func (e GameEngine) NewGame() (string, error) {
	resp, err := http.Get(config.CONFIG_SERVER_URL + "/game")
	if err != nil {
		log.Error(err)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return string(body), nil
}

func (e GameEngine) PlayTurn(gameId string, guessWord string) (string, error) {
	qs := fmt.Sprintf("/play?id=%s&guess=%s", gameId, guessWord)
	resp, err := http.Get(config.CONFIG_SERVER_URL + qs)
	if err != nil {
		log.Error(err)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return string(body), nil
}
