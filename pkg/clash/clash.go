package clash

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

type LockingCounter struct {
	sync.RWMutex
	Count map[string]int
}

//go:generate counterfeiter . Client
type Client interface {
	GetClan() (Clan, error)
	GetWar() (CurrentWar, error)
	CheckForAttackUpdates(*ClanWarMember, *LockingCounter) string
	GetOpponentMapPosition(tag string) int
}

type client struct {
	apiToken string
	clanTag  string
	baseURL  string
}

func NewClient(clanTag, apiToken, baseURL string) Client {
	return &client{
		clanTag:  clanTag,
		apiToken: apiToken,
		baseURL:  baseURL,
	}
}

func (c *client) GetClan() (Clan, error) {
	var clan Clan

	client := http.Client{}

	req, err := http.NewRequest("GET", c.baseURL+"/clans/"+url.QueryEscape(c.clanTag), nil)
	if err != nil {
		return clan, err
	}
	req.Header.Set("authorization", c.apiToken)
	resp, err := client.Do(req)
	if err != nil {
		return clan, err
	}

	if resp.StatusCode != http.StatusOK {
		return clan, errors.New("there was a problem getting the clan" + resp.Status)
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&clan)

	return clan, nil
}

func (c *client) GetWar() (CurrentWar, error) {
	var war CurrentWar

	client := http.Client{}

	req, err := http.NewRequest("GET", c.baseURL+"/clans/"+url.QueryEscape(c.clanTag)+"/currentwar", nil)
	if err != nil {
		return war, err
	}
	req.Header.Set("authorization", c.apiToken)
	resp, err := client.Do(req)
	if err != nil {
		return war, err
	}

	if resp.StatusCode != http.StatusOK {
		return war, errors.New("there was a problem getting the war" + resp.Status)
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&war)

	return war, nil
}

func (c *client) CheckForAttackUpdates(m *ClanWarMember, prevAttackCounter *LockingCounter) string {
	prevAttackCounter.RLock()
	newAttack := len(m.Attacks) > prevAttackCounter.Count[m.Tag]
	prevAttackCounter.RUnlock()

	prevAttackCounter.Lock()
	defer prevAttackCounter.Unlock()
	if newAttack {
		recentAttack := GetMostRecentAttack(m)
		defenderMapPosition := c.GetOpponentMapPosition(recentAttack.DefenderTag)

		prevAttackCounter.Count[m.Tag] = len(m.Attacks)
		return fmt.Sprintf("%s just %d starred number %d!\n", m.Name, recentAttack.Stars, defenderMapPosition)
	}

	prevAttackCounter.Count[m.Tag] = len(m.Attacks)
	return ""
}

func GetMostRecentAttack(m *ClanWarMember) ClanWarAttack {
	var recentAttack ClanWarAttack
	for _, a := range m.Attacks {
		if a.Order > recentAttack.Order {
			recentAttack = a
		}
	}
	return recentAttack
}

func (c *client) GetOpponentMapPosition(tag string) int {
	war, _ := c.GetWar()

	for _, p := range war.Opponent.Members {
		if p.Tag == tag {
			return p.MapPosition
		}
	}

	return -1
}

func GetWarResults(war *CurrentWar) string {
	if war.Clan.Stars > war.Opponent.Stars {
		return fmt.Sprintf("We won the war! The final star count was %d - %d", war.Clan.Stars, war.Opponent.Stars)
	} else if war.Clan.Stars < war.Opponent.Stars {
		return fmt.Sprintf("We lost the war ☹️.  The final star count was %d - %d", war.Clan.Stars, war.Opponent.Stars)
	} else {
		return fmt.Sprintf("We tied this war. The final star count was %d - %d", war.Clan.Stars, war.Opponent.Stars)
	}
}

func GetRemainingAttacks(war CurrentWar) map[string]int {
	attackMap := map[string]int{}

	for _, dudeOrDudette := range war.Clan.Members {
		attackMap[dudeOrDudette.Name] = 2 - len(dudeOrDudette.Attacks)
	}

	return attackMap
}
