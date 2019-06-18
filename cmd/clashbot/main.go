package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/robfig/cron"

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

	var state string
	var attackCounter = &clash.LockingCounter{Count: make(map[string]int)}

	c := cron.New()
	c.AddFunc("@every 30s", func() { state, attackCounter = RunBotLogic(clashClient, chatBot, state, attackCounter) })
	go c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

func RunBotLogic(clashClient clash.Client, chatBot chat.Bot, prevState string, prevAttackCounter *clash.LockingCounter) (string, *clash.LockingCounter) {
	war, err := clashClient.GetWar()
	if err != nil {
		fmt.Println("Error getting war", err)
	}

	if war.State == "inWar" {
		for _, m := range war.Clan.Members {
			go func(mem clash.ClanWarMember) {
				msg := clashClient.CheckForAttackUpdates(&mem, prevAttackCounter)
				fmt.Println("Attack updates:", msg)
				if msg != "" {
					chatBot.SendMessage(msg)
				}
			}(m)
		}
	}

	if war.State == "inWar" && prevState == "preparation" {
		chatBot.SendMessage("War has started!")
	}

	return war.State, prevAttackCounter
}
