package main

import (
	"fmt"
	"log"
	"os"

	"github.com/johnverrone/clashbot/bot"
	"github.com/johnverrone/clashbot/clash"
)

var BaseURL string

func main() {

	clanTag := os.Getenv("CLAN_TAG")
	clashToken := os.Getenv("CLASH_TOKEN")

	if clanTag == "" || clashToken == "" {
		log.Fatal("Environment not set correctly")
	}

	clashClient := clash.NewClashClient(clanTag, "Bearer "+clashToken, "https://api.clashofclans.com/v1")

	// clan, err := clashClient.GetClan()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	war, err := clashClient.GetWar()
	if err != nil {
		fmt.Println(err)
	}

	clashBot := bot.NewBot("groupme")
	fmt.Println("member:", war.Clan.Members[0].Name)
	fmt.Println("war:", war.Clan.Members[0].Attacks)

	msg, err := formatWarMessage(war)
	clashBot.SendMessage(msg)
}

func formatWarMessage(war clash.CurrentWar) (string, error) {

	return war.Clan.Members[0].Name, nil
}
