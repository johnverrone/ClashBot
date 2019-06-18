package chat

import (
	"log"
	"os"
)

type Bot interface {
	SendMessage(string) error
}

func NewBot(chatType string) Bot {
	botID := os.Getenv("BOT_ID")
	if botID == "" {
		log.Fatal("BOT_ID not set correctly")
	}
	switch chatType {
	case "groupme":
		return &GroupMeBot{
			token: "",
			botID: botID,
		}
	}
	return nil
}
