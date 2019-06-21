package clash

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

const (
	NotInWar    = "notInWar"
	Preparation = "preparation"
	InWar       = "inWar"
	WarEnded    = "warEnded"
)

type LockingCounter struct {
	sync.RWMutex
	Count map[string]int
}

type Location struct {
	id        int
	name      string
	isCountry bool
}

type WarLog struct {
	Items []string `json:"items"`
}

type CurrentWar struct {
	State                string  `json:"state"`
	TeamSize             int     `json:"teamSize"`
	PreparationStartTime string  `json:"preparationStartTime"`
	StartTime            string  `json:"startTime"`
	EndTime              string  `json:"endTime"`
	Clan                 WarClan `json:"clan"`
	Opponent             WarClan `json:"opponent"`
}

type WarClan struct {
	Tag       string          `json:"tag"`
	Name      string          `json:"name"`
	BadgeUrls URLContainer    `json:"badgeUrls"`
	ClanLevel int             `json:"clanLevel"`
	Attacks   int             `json:"attacks"`
	Stars     int             `json:"stars"`
	ExpEarned int             `json:"expEarned"`
	Members   []ClanWarMember `json:"members"`
}

type URLContainer struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

type ClanWarMember struct {
	Tag                string          `json:"tag"`
	Name               string          `json:"name"`
	TownHallLevel      int             `json:"townHallLevel"`
	MapPosition        int             `json:"mapPosition"`
	Attacks            []ClanWarAttack `json:"attacks"`
	OpponenetAttacks   int             `json:"opponenetAttacks"`
	BestOpponentAttack ClanWarAttack   `json:"bestOpponentAttack"`
}

type ClanWarAttack struct {
	AttackerTag           string `json:"attackerTag"`
	DefenderTag           string `json:"defenderTag"`
	Stars                 int    `json:"stars"`
	DestructionPercentage int    `json:"destructionPercentage"`
	Order                 int    `json:"order"`
}

type Clan struct {
	Tag              string   `json:"tag"`
	Name             string   `json:"name"`
	Location         Location `json:"location"`
	ClanLevel        int      `json:"clanLevel"`
	ClanPoints       int      `json:"clanPoints"`
	ClanVersusPoints int      `json:"clanVersusPoints"`
	Members          int      `json:"members"`
	ClanType         string   `json:"clanType"`
	RequiredTrophies int      `json:"requiredTrophies"`
	WarFrequency     string   `json:"warFrequency"`
	WarWinStreak     int      `json:"warWinStreak"`
	WarWins          int      `json:"warWins"`
	WarTies          int      `json:"warTies"`
	WarLosses        int      `json:"warLosses"`
	IsWarLogPublic   bool     `json:"isWarLogPublic"`
	Description      string   `json:"description"`
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
	newAttack := len(m.Attacks) > prevAttackCounter.Count[m.Name]
	prevAttackCounter.RUnlock()

	prevAttackCounter.Lock()
	defer prevAttackCounter.Unlock()
	if newAttack {
		recentAttack := GetMostRecentAttack(m)
		defenderMapPosition := c.GetOpponentMapPosition(recentAttack.DefenderTag)

		prevAttackCounter.Count[m.Name] = len(m.Attacks)
		return fmt.Sprintf("%s just %d starred number %d!\n", m.Name, recentAttack.Stars, defenderMapPosition)
	}

	prevAttackCounter.Count[m.Name] = len(m.Attacks)
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
