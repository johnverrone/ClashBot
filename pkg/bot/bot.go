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

func RunBotLogic(clashClient clash.Client, chatClient chat.Client, prevState *PrevState) {
	war, err := clashClient.GetWar()
	if err != nil {
		fmt.Println("Error getting war", err)
	}

	if war.State == clash.NotInWar && prevState.War == clash.InWar {
		msg := clashClient.GetWarResults(&war)
		if msg != "" {
			chatClient.SendMessage(msg)
		}
		return
	}

	if war.State == clash.InWar {
		for _, m := range war.Clan.Members {
			go func(mem clash.ClanWarMember) {
				msg := clashClient.CheckForAttackUpdates(&mem, prevState.AttackCounter)
				fmt.Println("Attack updates:", msg)
				if msg != "" {
					chatClient.SendMessage(msg)
				}
			}(m)
		}
	}

	if war.State == clash.InWar && prevState.War == clash.Preparation {
		chatClient.SendMessage("War has started!")
	}

	prevState.War = war.State
}
