package store

import (
	"testing"

	"aluance.io/wordleplayer/internal/playerbot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitStore(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	require.NoError(useTestMode())
	cleanupTestMode()

	sm, err := createOneGameSM()
	require.NoError(err)

	v, ok := sm.(*oneGameSM)
	require.True(ok)

	err = v.initStore()
	assert.NoError(err)
}

func TestSave(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	require.NoError(useTestMode())
	// cleanupTestMode()

	sm, err := createOneGameSM()
	require.NoError(err)

	type oneTurnTEST struct {
		Guess     string   `dynamodbav:"guess"`
		IsValid   bool     `dynamodbav:"isValid"`
		TryResult []string `dynamodbav:"tryResult"`
	}

	type oneGameTEST struct {
		PlayerName     string        `dynamodbav:"playerName"`
		GameId         string        `dynamodbav:"gameId"`
		GameStatus     string        `dynamodbav:"gameStatus"`
		Turns          []oneTurnTEST `dynamodbav:"turns"`
		WinWord        string        `dynamodbav:"winWord"`
		ValidAttempts  int           `dynamodbav:"validAttempts"`
		WinningAttempt int           `dynamodbav:"winningAttempt"`
	}

	testData := oneGameTEST{
		PlayerName:    playerbot.ONEBOT_NAME,
		GameId:        "id001_TestOneGameSmSave",
		GameStatus:    "Lost",
		WinWord:       "BLAHS",
		ValidAttempts: 6,
	}
	err = sm.Save(&testData)
	assert.NoError(err)
}
