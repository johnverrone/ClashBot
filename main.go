package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/robfig/cron"

	"github.com/johnverrone/clashbot/pkg/bot"
	"github.com/johnverrone/clashbot/pkg/chat"
	"github.com/johnverrone/clashbot/pkg/clash"
)

func main() {
	clanTag := os.Getenv("CLAN_TAG")
	clashAPIKey := os.Getenv("CLASH_API_KEY")

	if clanTag == "" || clashAPIKey == "" {
		log.Fatal("CLAN_TAG or CLASH_API_KEY not set correctly")
	}

	clashClient := clash.NewClient(clanTag, "Bearer "+clashAPIKey, "https://api.clashofclans.com/v1")
	chatClient := chat.NewClient("groupme")

	prevState := bot.PrevState{
		War:           "",
		AttackCounter: &clash.LockingCounter{Count: make(map[string]int)},
	}

	c := cron.New()
	c.AddFunc("@every 30s", func() { bot.RunBotLogic(clashClient, chatClient, &prevState) })
	go c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
