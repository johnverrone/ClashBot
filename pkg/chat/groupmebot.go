package chat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type GroupMeBot struct {
	token string
	botID string
}

type PostBody struct {
	BotID string `json:"bot_id"`
	Text  string `json:"text"`
}

func (b *GroupMeBot) SendMessage(message string) error {

	postBody := PostBody{
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
