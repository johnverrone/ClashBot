package bot

import (
	"fmt"

	"github.com/johnverrone/clashbot/pkg/chat"
	"github.com/johnverrone/clashbot/pkg/clash"
)

type PrevState struct {
	War           string
	AttackCounter *clash.LockingCounter
}

func RunBotLogic(clashClient clash.Client, chatBot chat.Bot, prevState *PrevState) {
	war, err := clashClient.GetWar()
	if err != nil {
		fmt.Println("Error getting war", err)
	}

	if war.State == "notInWar" && prevState.War == "inWar" {
		msg := clashClient.GetWarResults(&war)
		if msg != "" {
			chatBot.SendMessage(msg)
		}
		return
	}

	if war.State == "inWar" {
		for _, m := range war.Clan.Members {
			go func(mem clash.ClanWarMember) {
				msg := clashClient.CheckForAttackUpdates(&mem, prevState.AttackCounter)
				fmt.Println("Attack updates:", msg)
				if msg != "" {
					chatBot.SendMessage(msg)
				}
			}(m)
		}
	}

	if war.State == "inWar" && prevState.War == "preparation" {
		chatBot.SendMessage("War has started!")
	}

	prevState.War = war.State
}
