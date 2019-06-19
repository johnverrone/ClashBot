package chat

import (
	"log"
	"os"
)

type Client interface {
	SendMessage(string) error
}

func NewClient(chatType string) Client {
	switch chatType {
	case "groupme":
		botID := os.Getenv("GROUPME_BOT_ID")
		if botID == "" {
			log.Fatal("GROUPME_BOT_ID not set correctly")
		}
		return NewGroupMeClient(botID)
	}
	return nil
}
