package playerbot

import (
	"testing"

	"aluance.io/wordleplayer/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateDictionary(t *testing.T) {
	assert := assert.New(t)

	startSize := len(dictionaries)

	dummyBotId := "a1b2c3d4e5_TestCreateDictionary"
	d := CreateDictionary(dummyBotId)
	assert.Implements((*PlayerDictionary)(nil), d)
	assert.Equal(startSize+1, len(dictionaries))

	dummyWord := "ABCDE"
	d.Remember(dummyWord, true)

	d2 := CreateDictionary(dummyBotId)
	assert.Equal(startSize+1, len(dictionaries))
	assert.True(d2.IsValid(dummyWord))
}

func TestRemember(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		word  string
		valid bool
		err   error
	}{
		{word: "ABCDE", valid: false, err: nil},
		{word: "WHALE", valid: true, err: nil},
	}

	dummyBotId := "a1b2c3d4e5_TestRemember"
	d := CreateDictionary(dummyBotId)

	for _, test := range tests {
		err := d.Remember(test.word, test.valid)
		assert.NoError(err)

		v, ok := d.(*playerDictionary)
		require.True(ok)

		valid, ok := v.lexicon[test.word]
		assert.True(ok)
		assert.Equal(test.valid, valid)
	}
}

func TestGenerate(t *testing.T) {
	assert := assert.New(t)

	dummyBotId := "a1b2c3d4e5_TestGenerate"
	d := CreateDictionary(dummyBotId)

	w, err := d.Generate()
	assert.NoError(err)
	assert.Equal(config.CONFIG_GAME_WORDLENGTH, len(w))

	w2, err := d.Generate()
	assert.NoError(err)
	assert.Equal(config.CONFIG_GAME_WORDLENGTH, len(w2))
	assert.NotEqualValues(w, w2)
}

func TestIsValid(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		word  string
		valid bool
	}{
		{word: "ABCDE", valid: false},
		{word: "WHALE", valid: true},
	}

	dummyBotId := "a1b2c3d4e5_TestIsValid"
	d := CreateDictionary(dummyBotId)

	// Add words to dictionary
	for _, test := range tests {
		err := d.Remember(test.word, test.valid)
		require.NoError(err)
	}

	// Test that the word validity is correct
	for _, test := range tests {
		valid := d.IsValid(test.word)
		assert.Equal(test.valid, valid)
	}

	// Test an unknown word
	assert.False(d.IsValid("WRONG"))
}
