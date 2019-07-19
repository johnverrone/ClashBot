package chat

import (
	"log"
	"net/http"
	"os"
)

//go:generate counterfeiter . Client
type Client interface {
	SendMessage(string) error
	HandleMessage(http.ResponseWriter, *http.Request)
}

func NewClient(chatType string) Client {
	switch chatType {
	case "groupme":
		botID := os.Getenv("GROUPME_BOT_ID")
		if botID == "" {
			log.Fatal("GROUPME_BOT_ID not set correctly")
		}
		return NewGroupMeClient(botID)
	case "slack":
		webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
		if webhookURL == "" {
			log.Fatal("SLACK_WEBHOOK_URL not set correctly")
		}
		return NewSlackClient(webhookURL)
	}
	return nil
}
