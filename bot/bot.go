package bot

import (
	"log"
	"os"
)

type Bot interface {
	SendMessage(string) error
}

func NewBot(chatType string) Bot {
	botId := os.Getenv("BOT_ID")
	if botId == "" {
		log.Fatal("Environment not set correctly")
	}
	switch chatType {
	case "groupme":
		return &GroupMeBot{
			token: "",
			botId: botId,
		}
	}
	return nil
}
