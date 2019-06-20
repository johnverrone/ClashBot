package chat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type SlackClient struct {
	webhookURL string
}

type SlackPostBody struct {
	Text string `json:"text"`
}

func NewSlackClient(webhookURL string) *SlackClient {
	return &SlackClient{
		webhookURL: webhookURL,
	}
}

func (c *SlackClient) SendMessage(message string) error {

	postBody := SlackPostBody{
		Text: message,
	}

	data, err := json.Marshal(postBody)
	if err != nil {
		fmt.Println(err)
		return err
	}

	resp, err := http.Post(c.webhookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		return err
	}

	if resp.StatusCode/100 != 2 {
		fmt.Println("ERROR:", resp)
		return errors.New("there was a problem posting to slack" + resp.Status)
	}

	return nil
}
