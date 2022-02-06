package playerbot

import "aluance.io/wordleplayer/internal/config"

type Playerbot interface {
	PlayGame(ch *chan string)
}

type playerbotFactory func() (Playerbot, error)

var mapBots = map[string]playerbotFactory{
	ONEBOT_NAME: createOne,
	TWOBOT_NAME: createTwo,
}

func CreateBot(name string, options ...interface{}) (Playerbot, error) {
	bf, ok := mapBots[name]
	if !ok {
		return nil, ErrInvalidBotName
	}

	bot, err := bf()
	if err != nil {
		return nil, err
	}

	return bot, nil
}

type botGame struct {
	PlayerName     string    `dynamodbav:"playerName"`
	GameId         string    `dynamodbav:"gameId"`
	GameStatus     string    `dynamodbav:"gameStatus"`
	Turns          []botTurn `dynamodbav:"turns"`
	WinWord        string    `dynamodbav:"winWord"`
	ValidAttempts  int       `dynamodbav:"validAttempts"`
	WinningAttempt int       `dynamodbav:"winningAttempt"`
}

type botTurn struct {
	Guess     string   `dynamodbav:"guess"`
	IsValid   bool     `dynamodbav:"isValid"`
	TryResult []string `dynamodbav:"tryResult"`
}

func createBotTurn() *botTurn {
	turn := new(botTurn)
	turn.TryResult = make([]string, config.CONFIG_GAME_WORDLENGTH)

	return turn
}
