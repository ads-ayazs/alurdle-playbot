package playerbot

type Playerbot interface {
	PlayGame(ch *chan string)
}

type playerbotFactory func() (Playerbot, error)

var mapBots = map[string]playerbotFactory{
	ONEBOT_NAME: createOne,
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
