package xplorersbot

import (
	"testing"
)

// type SlackStruct struct{}

// func (sl SlackStruct) New(token string) *slack.Client {
// 	return &slack.Client{}
// }

func TestGetSlackApiClient(t *testing.T) {
	// assert := require.New(t)
	// sl := SlackStruct{}
	cmd, err := GetSlackApiClient()
	if err != nil {
		t.Log("Unable to fetch slack client")
		t.Log(err)
		return
	}
	t.Log(cmd)
	// assert.Equal(*result.Parameter.Value, "some-parameter-value")
}
