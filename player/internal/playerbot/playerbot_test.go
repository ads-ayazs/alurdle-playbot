package playerbot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//func CreateBot(name string, options ...interface{}) (Playerbot, error) {

func TestCreateBot(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name    string
		options []interface{}
		result  interface{}
		err     error
	}{
		{name: ONEBOT_NAME, result: (*Playerbot)(nil), err: nil},
	}

	for _, test := range tests {
		pb, err := CreateBot(test.name, test.options...)
		if test.err != nil {
			assert.ErrorIs(err, test.err)
			continue
		}
		assert.NoError(err)
		assert.Implements(test.result, pb)
	}
}
