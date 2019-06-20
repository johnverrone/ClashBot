package chat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type GroupMeClient struct {
	botID string
}

type GroupMePostBody struct {
	BotID string `json:"bot_id"`
	Text  string `json:"text"`
}

func NewGroupMeClient(botID string) *GroupMeClient {
	return &GroupMeClient{
		botID: botID,
	}
}

func (b *GroupMeClient) SendMessage(message string) error {

	postBody := GroupMePostBody{
		BotID: b.botID,
		Text:  message,
	}

	data, err := json.Marshal(postBody)
	if err != nil {
		fmt.Println(err)
		return err
	}

	resp, err := http.Post("https://api.groupme.com/v3/bots/post", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		return err
	}

	if resp.StatusCode/100 != 2 {
		fmt.Println("ERROR:", resp)
		return errors.New("there was a problem posting to groupme" + resp.Status)
	}

	return nil
}
