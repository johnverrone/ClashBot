package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

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

	// initialize War Status and Attack Counter when starting the bot to avoid duplicate messages
	war, _ := clashClient.GetWar()
	for _, m := range war.Clan.Members {
		prevState.AttackCounter.Count[m.Name] = len(m.Attacks)
	}

	fmt.Printf("Bot has started at %s.\nThe current war status is: %s\nCurrent attack count: %d - %d\nCurrent star count: %d - %d\n\n", time.Now(), war.State, war.Clan.Attacks, war.Opponent.Attacks, war.Clan.Stars, war.Opponent.Stars)
	fmt.Println("Current attack count breakdown:")
	for pl, count := range prevState.AttackCounter.Count {
		fmt.Printf("%s: %d\n", pl, count)
	}

	c := cron.New()
	c.AddFunc("@every 30s", func() { bot.RunBotLogic(clashClient, chatClient, &prevState) })
	go c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
