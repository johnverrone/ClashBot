package clash

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

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
	BadgeUrls UrlContainer    `json:"badgeUrls"`
	ClanLevel int             `json:"clanLevel"`
	Attacks   int             `json:"endTime"`
	Stars     int             `json:"clan"`
	ExpEarned int             `json:"expEarned"`
	Members   []ClanWarMember `json:"members"`
}

type UrlContainer struct {
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

type Clash struct {
	apiToken string
	clanTag  string
	baseUrl  string
}

func NewClashClient(clanTag, apiToken, baseUrl string) Clash {
	return Clash{
		clanTag:  clanTag,
		apiToken: apiToken,
		baseUrl:  baseUrl,
	}
}

func (c *Clash) GetClan() (Clan, error) {
	var clan Clan

	client := http.Client{}

	req, err := http.NewRequest("GET", c.baseUrl+"/clans/"+url.QueryEscape(c.clanTag), nil)
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

func (c *Clash) GetWar() (CurrentWar, error) {
	var war CurrentWar

	client := http.Client{}

	req, err := http.NewRequest("GET", c.baseUrl+"/clans/"+url.QueryEscape(c.clanTag)+"/currentwar", nil)
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

func (c *Clash) CheckForWar(state chan<- string) {

	war, err := c.GetWar()
	if err != nil {
		fmt.Println("Error checking war status: ", err)
	}

	state <- war.State
}
