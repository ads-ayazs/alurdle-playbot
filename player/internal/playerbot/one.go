package playerbot

import (
	"encoding/json"

	"aluance.io/wordleplayer/internal/config"
	"aluance.io/wordleplayer/internal/store"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
)

const ONEBOT_NAME = "one"

type oneBot struct {
	Id         string
	Game       oneGame
	Dictionary PlayerDictionary
}

func createOne() (Playerbot, error) {
	bot := new(oneBot)
	bot.Id = xid.New().String()
	bot.Game.PlayerName = ONEBOT_NAME
	log.Info("botId: ", bot.Id)

	bot.Dictionary = CreateDictionary(bot.Id)
	if bot.Dictionary == nil {
		return nil, ErrNilDictionary
	}

	return bot, nil
}

func (b oneBot) PlayGame(ch *chan string) {
	// Avoid sending to nil channel
	if ch == nil {
		return
	}

	// Start the game
	if err := b.startGame(); err != nil {
		*ch <- ""
		return
	}

	// Play while the game remains "InPlay"
	for b.isGameInPlay() {
		if err := b.playTurn(); err != nil {
			*ch <- ""
			return
		}
	}

	// Finish the game and write the output to ch
	*ch <- b.finishGame()
}

type oneGame struct {
	PlayerName     string    `dynamodbav:"playerName"`
	GameId         string    `dynamodbav:"gameId"`
	GameStatus     string    `dynamodbav:"gameStatus"`
	Turns          []oneTurn `dynamodbav:"turns"`
	WinWord        string    `dynamodbav:"winWord"`
	ValidAttempts  int       `dynamodbav:"validAttempts"`
	WinningAttempt int       `dynamodbav:"winningAttempt"`
}

type oneTurn struct {
	Guess     string   `dynamodbav:"guess"`
	IsValid   bool     `dynamodbav:"isValid"`
	TryResult []string `dynamodbav:"tryResult"`
}

func createOneTurn() *oneTurn {
	turn := new(oneTurn)
	turn.TryResult = make([]string, config.CONFIG_GAME_WORDLENGTH)

	return turn
}

func (bot *oneBot) startGame() error {
	ge := GetGameEngine()

	// Create a new game and save the game id
	out, err := ge.NewGame()
	if err != nil {
		return err
	}

	// Unmarshall the output
	outmap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(out), &outmap); err != nil {
		return err
	}

	// Save essential information
	bot.Game.GameId = outmap["id"].(string)
	bot.Game.GameStatus = outmap["gameStatus"].(string)

	return nil
}

func (bot oneBot) isGameInPlay() bool {
	return bot.Game.GameStatus == "InPlay"
}

func (bot *oneBot) playTurn() error {
	ge := GetGameEngine()

	// Generate a word
	if bot.Dictionary == nil {
		return ErrNilDictionary
	}
	guessWord, err := bot.Dictionary.Generate()
	if err != nil {
		return err
	}

	// Play the guess word
	out, err := ge.PlayTurn(bot.Game.GameId, guessWord)
	if err != nil {
		return err
	}

	// Unmarshall the output
	outmap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(out), &outmap); err != nil {
		return err
	}

	// Find the latest attempt
	attempts := outmap["attempts"].([]interface{})
	lastAttempt := attempts[len(attempts)-1].(map[string]interface{})
	tw := lastAttempt["tryWord"].(string)
	if tw != guessWord {
		return ErrFailedAttempt
	}

	// Create and save a turn record
	turn := createOneTurn()
	turn.Guess = guessWord
	turn.IsValid = lastAttempt["isValidWord"].(bool)

	tr := lastAttempt["tryResult"].([]interface{})
	for i := 0; i < len(tr); i++ {
		turn.TryResult[i] = tr[i].(string)
	}

	bot.Game.Turns = append(bot.Game.Turns, *turn)

	// Update the dictionary
	if err := bot.Dictionary.Remember(turn.Guess, turn.IsValid); err != nil {
		return err
	}

	// Save essential information
	bot.Game.GameStatus = outmap["gameStatus"].(string)
	bot.Game.ValidAttempts = int(outmap["validAttempts"].(float64))
	if bot.Game.GameStatus == "Won" {
		bot.Game.WinWord = outmap["secretWord"].(string)
		bot.Game.WinningAttempt = outmap["winningAttempt"].(int)
	} else if bot.Game.GameStatus == "Lost" {
		bot.Game.WinWord = outmap["secretWord"].(string)

		// Update the dictionary
		if err := bot.Dictionary.Remember(bot.Game.WinWord, true); err != nil {
			return err
		}
	}

	return nil
}

func (bot oneBot) finishGame() string {

	sm, err := store.GetStoreManager(ONEBOT_NAME)
	if err == nil {
		if err := sm.Save(bot.Game); err != nil {
			log.Error(err)
		}
	} else {
		log.Error(err)
	}

	log.Info("BOT FINISHED - ", "botId: ", bot.Id, " gameId: ", bot.Game.GameId)
	log.Info("    botId: ", bot.Id, " dictionary valid/size: ", bot.Dictionary.DescribeSize(true), "/", bot.Dictionary.DescribeSize(false))
	log.Info("    botId: ", bot.Id, " outcome: ", bot.Game.GameStatus)
	return bot.Game.GameId
}
