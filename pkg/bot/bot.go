package bot

import (
	"fmt"
	"time"

	"github.com/johnverrone/clashbot/pkg/chat"
	"github.com/johnverrone/clashbot/pkg/clash"
)

type PrevState struct {
	War           string
	AttackCounter *clash.LockingCounter
}

func RunBotLogic(clashClient clash.Client, chatClient chat.Client, prevState *PrevState) {
	war, err := clashClient.GetWar()
	if err != nil {
		fmt.Println("Error getting war", err)
	}

	defer func() { prevState.War = war.State }()

	switch true {
	case war.State == clash.NotInWar && prevState.War == clash.InWar:
		handleWarEnd(clashClient, chatClient, war)
		return
	case war.State == clash.NotInWar:
		// do nothing...
		return
	case war.State == clash.InWar && prevState.War == clash.Preparation:
		handleWarStart(chatClient)
	case war.State == clash.InWar:
		handleWarUpdates(clashClient, chatClient, prevState, war)
		return
	}
}

func handleWarStart(chatClient chat.Client) {
	chatClient.SendMessage("War has started!")
}

func handleWarEnd(clashClient clash.Client, chatClient chat.Client, war clash.CurrentWar) {
	msg := clashClient.GetWarResults(&war)
	if msg != "" {
		chatClient.SendMessage(msg)
	}
}

func handleWarUpdates(clashClient clash.Client, chatClient chat.Client, prevState *PrevState, war clash.CurrentWar) {
	for _, m := range war.Clan.Members {
		go func(mem clash.ClanWarMember) {
			msg := clashClient.CheckForAttackUpdates(&mem, prevState.AttackCounter)
			if msg != "" {
				fmt.Printf("%s: %s\n", time.Now(), msg)
				chatClient.SendMessage(msg)
			}
		}(m)
	}

	layout := "20060102T150405.000Z"
	t, err := time.Parse(layout, war.EndTime)
	if err != nil {
		fmt.Println("Error parsing time", err)
	}
	twoHours, _ := time.ParseDuration("2h")
	if time.Until(t) < twoHours {
		fmt.Printf("Less than 2 hours remaining: %v\n", time.Until(t))
	}
}
