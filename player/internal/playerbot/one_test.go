package playerbot

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateOne(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		result oneBot
		err    error
	}{
		{result: oneBot{}, err: nil},
	}

	for _, test := range tests {
		pb, err := createOne()
		assert.NoError(err)
		assert.IsType(&test.result, pb)

		v, ok := pb.(*oneBot)
		assert.True(ok)
		assert.NotEmpty(v.id)
	}
}

func TestOneBotPlayGame(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	pb, err := createOne()
	require.NoError(err)

	ch := make(chan string)
	go pb.PlayGame(&ch)
	defer close(ch)

	select {
	case s := <-ch:
		assert.NotEmpty(s)
	case <-time.After(1000 * time.Second):
		assert.Fail("timed out without receiving from channel")
	}
}

func TestOneBotStartGame(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	pb, err := createOne()
	require.NoError(err)

	v, ok := pb.(*oneBot)
	require.True(ok)

	err = v.startGame()
	assert.NoError(err)

	assert.NotEmpty(v.game.gameId)
	assert.Equal("InPlay", v.game.gameStatus)
}

func TestOneBotIsGameInPlay(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	pb, err := createOne()
	require.NoError(err)

	v, ok := pb.(*oneBot)
	require.True(ok)

	inPlay := v.isGameInPlay()
	assert.Equal(inPlay, v.game.gameStatus == "InPlay")
}

func TestOneBotPlayTurn(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	pb, err := createOne()
	require.NoError(err)

	v, ok := pb.(*oneBot)
	require.True(ok)

	err = v.startGame()
	require.NoError(err)

	err = v.playTurn()
	assert.NoError(err)

	assert.NotZero(len(v.game.turns))
}

func TestOneBotFinishGame(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	pb, err := createOne()
	require.NoError(err)

	v, ok := pb.(*oneBot)
	require.True(ok)

	err = v.startGame()
	require.NoError(err)

	err = v.playTurn()
	require.NoError(err)

	s := v.finishGame()
	assert.Equal(v.game.gameId, s)
}
