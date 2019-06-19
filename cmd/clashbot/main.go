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
	clashToken := os.Getenv("CLASH_TOKEN")

	if clanTag == "" || clashToken == "" {
		log.Fatal("CLAN_TAG or CLASH_TOKEN not set correctly")
	}

	clashClient := clash.NewClient(clanTag, "Bearer "+clashToken, "https://api.clashofclans.com/v1")
	chatBot := chat.NewBot("groupme")

	prevState := bot.PrevState{
		War:           "",
		AttackCounter: &clash.LockingCounter{Count: make(map[string]int)},
	}

	c := cron.New()
	c.AddFunc("@every 30s", func() { bot.RunBotLogic(clashClient, chatBot, &prevState) })
	go c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
