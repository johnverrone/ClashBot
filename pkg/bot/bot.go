package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/johnverrone/clashbot/pkg/chat"
	"github.com/johnverrone/clashbot/pkg/clash"
)

type PrevState struct {
	War             string
	AttackCounter   *clash.LockingCounter
	SentWarReminder bool
}

func RunBotLogic(clashClient clash.Client, chatClient chat.Client, prevState *PrevState) error {
	war, err := clashClient.GetWar()
	if err != nil {
		log.Println("Error getting war:", err)
		return err
	}

	log.Println("War status is", war.State)

	defer func() { prevState.War = war.State }()

	switch true {
	case war.State == clash.WarEnded && prevState.War == clash.InWar:
		handleWarEnd(clashClient, chatClient, war)
	case war.State == clash.WarEnded || war.State == clash.NotInWar:
		// do nothing...
	case war.State == clash.InWar && prevState.War == clash.Preparation:
		handleWarStart(chatClient)
	case war.State == clash.InWar:
		handleWarUpdates(clashClient, chatClient, prevState, war)
	}
	return nil
}

func handleWarStart(chatClient chat.Client) {
	chatClient.SendMessage("War has started!")
}

func handleWarEnd(clashClient clash.Client, chatClient chat.Client, war clash.CurrentWar) {
	msg := clash.GetWarResults(&war)
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

	sendEndOfWarReminder(prevState, chatClient, war)
}

func sendEndOfWarReminder(prevState *PrevState, chatClient chat.Client, war clash.CurrentWar) {
	twoHours, _ := time.ParseDuration("2h")
	remainingTime, err := getRemainingWarTime(war)
	if err != nil {
		return
	}

	if !prevState.SentWarReminder && remainingTime < twoHours {
		remainingAttacks := `There is less than 2 hours remaining in the war.`
		attackMap := clash.GetRemainingAttacks(war)
		for member, numAttacks := range attackMap {
			if numAttacks > 0 {
				remainingAttacks += `\n` + member
			}
		}

		chatClient.SendMessage(remainingAttacks)
		prevState.SentWarReminder = true
	}
}

func getRemainingWarTime(war clash.CurrentWar) (time.Duration, error) {
	layout := "20060102T150405.000Z"
	t, err := time.Parse(layout, war.EndTime)
	if err != nil {
		fmt.Println("Error parsing time", err)
		return time.Minute, err
	}

	return time.Until(t), nil
}
