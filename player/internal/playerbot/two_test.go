package playerbot

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTwo(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		result twoBot
		err    error
	}{
		{result: twoBot{}, err: nil},
	}

	for _, test := range tests {
		pb, err := createTwo()
		assert.NoError(err)
		assert.IsType(&test.result, pb)

		v, ok := pb.(*twoBot)
		assert.True(ok)
		assert.NotEmpty(v.Id)
	}
}

func TestTwoBotPlayGame(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	pb, err := createTwo()
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

func TestTwoBotStartGame(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	pb, err := createTwo()
	require.NoError(err)

	v, ok := pb.(*twoBot)
	require.True(ok)

	err = v.startGame()
	assert.NoError(err)

	assert.NotEmpty(v.Game.GameId)
	assert.Equal("InPlay", v.Game.GameStatus)
}

func TestTwoBotIsGameInPlay(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	pb, err := createTwo()
	require.NoError(err)

	v, ok := pb.(*twoBot)
	require.True(ok)

	inPlay := v.isGameInPlay()
	assert.Equal(inPlay, v.Game.GameStatus == "InPlay")
}

func TestTwoBotPlayTurn(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	pb, err := createTwo()
	require.NoError(err)

	v, ok := pb.(*twoBot)
	require.True(ok)

	err = v.startGame()
	require.NoError(err)

	err = v.playTurn()
	assert.NoError(err)

	assert.NotZero(len(v.Game.Turns))
}

func TestTwoBotFinishGame(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	pb, err := createTwo()
	require.NoError(err)

	v, ok := pb.(*twoBot)
	require.True(ok)

	err = v.startGame()
	require.NoError(err)

	err = v.playTurn()
	require.NoError(err)

	s := v.finishGame()
	assert.Equal(v.Game.GameId, s)
}
