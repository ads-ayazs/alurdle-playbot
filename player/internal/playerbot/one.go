package playerbot

import (
	"encoding/json"

	"aluance.io/wordleplayer/internal/config"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
)

const ONEBOT_NAME = "one"

type oneBot struct {
	id         string
	game       oneGame
	dictionary PlayerDictionary
}

func createOne() (Playerbot, error) {
	bot := new(oneBot)
	bot.id = xid.New().String()
	log.Info("botId: ", bot.id)

	bot.dictionary = CreateDictionary(bot.id)
	if bot.dictionary == nil {
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
	gameId         string
	gameStatus     string
	turns          []oneTurn
	winWord        string
	validAttempts  int
	winningAttempt int
}

type oneTurn struct {
	guess     string
	isValid   bool
	tryResult []string
}

func createOneTurn() *oneTurn {
	turn := new(oneTurn)
	turn.tryResult = make([]string, config.CONFIG_GAME_WORDLENGTH)

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
	bot.game.gameId = outmap["id"].(string)
	bot.game.gameStatus = outmap["gameStatus"].(string)

	return nil
}

func (bot oneBot) isGameInPlay() bool {
	return bot.game.gameStatus == "InPlay"
}

func (bot *oneBot) playTurn() error {
	ge := GetGameEngine()

	// Generate a word
	if bot.dictionary == nil {
		return ErrNilDictionary
	}
	guessWord, err := bot.dictionary.Generate()
	if err != nil {
		return err
	}

	// Play the guess word
	out, err := ge.PlayTurn(bot.game.gameId, guessWord)
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
	turn.guess = guessWord
	turn.isValid = lastAttempt["isValidWord"].(bool)

	tr := lastAttempt["tryResult"].([]interface{})
	for i := 0; i < len(tr); i++ {
		turn.tryResult[i] = tr[i].(string)
	}

	bot.game.turns = append(bot.game.turns, *turn)

	// Update the dictionary
	if err := bot.dictionary.Remember(turn.guess, turn.isValid); err != nil {
		return err
	}

	// Save essential information
	bot.game.gameStatus = outmap["gameStatus"].(string)
	bot.game.validAttempts = int(outmap["validAttempts"].(float64))
	if bot.game.gameStatus == "Won" {
		bot.game.winWord = outmap["secretWord"].(string)
		bot.game.winningAttempt = outmap["winningAttempt"].(int)
	} else if bot.game.gameStatus == "Lost" {
		bot.game.winWord = outmap["secretWord"].(string)

		// Update the dictionary
		if err := bot.dictionary.Remember(bot.game.winWord, true); err != nil {
			return err
		}
	}

	return nil
}

func (bot oneBot) finishGame() string {

	log.Info("BOT FINISHED - ", "botId: ", bot.id, " gameId: ", bot.game.gameId)
	log.Info("    botId: ", bot.id, " dictionary valid/size: ", bot.dictionary.DescribeSize(true), "/", bot.dictionary.DescribeSize(false))
	log.Info("    botId: ", bot.id, " outcome: ", bot.game.gameStatus)
	return bot.game.gameId
}
