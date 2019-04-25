package main

import (
	"fmt"
	"log"
	"os"

	"github.com/johnverrone/clashbot/bot"
	"github.com/johnverrone/clashbot/clash"
	"github.com/robfig/cron"
)

var BaseURL string

func main() {

	clanTag := os.Getenv("CLAN_TAG")
	clashToken := os.Getenv("CLASH_TOKEN")

	if clanTag == "" || clashToken == "" {
		log.Fatal("Environment not set correctly")
	}

	clashClient := clash.NewClashClient(clanTag, "Bearer "+clashToken, "https://api.clashofclans.com/v1")

	warState := make(chan string)
	c := cron.New()
	c.AddFunc("*/10 * * * * *", func() { clashClient.CheckForWar(warState) })
	go c.Start()
	defer c.Stop()

	clashBot := bot.NewBot("groupme")

	var state string
	var prevState string
	for {
		prevState = state
		state = <-warState

		war, err := clashClient.GetWar()
		if err != nil {
			fmt.Println("Error getting war", err)
			continue
		}

		if state == "inWar" {

			var prevAttackCounter = &clash.LockingCounter{Count: make(map[string]int)}

			for _, m := range war.Clan.Members {
				prevAttackCounter.Count[m.Name] = len(m.Attacks)
			}

			for _, m := range war.Clan.Members {
				go func(mem clash.ClanWarMember) {
					msg := clashClient.CheckForAttackUpdates(&mem, prevAttackCounter)
					if msg != "" {
						clashBot.SendMessage(msg)
					}
				}(m)
			}
		}

		if prevState == "preparation" && state == "inWar" {
			clashBot.SendMessage("War has started!")
			break
		}
	}
}
