package clash

const (
	NotInWar    = "notInWar"
	Preparation = "preparation"
	InWar       = "inWar"
	WarEnded    = "warEnded"
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
