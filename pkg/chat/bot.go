package chat

import (
	"log"
	"os"
)

type Bot interface {
	SendMessage(string) error
}

func NewBot(chatType string) Bot {
	switch chatType {
	case "groupme":
		botID := os.Getenv("GROUPME_BOT_ID")
		if botID == "" {
			log.Fatal("GROUPME_BOT_ID not set correctly")
		}
		return &GroupMeBot{
			token: "",
			botID: botID,
		}
	}
	return nil
}
