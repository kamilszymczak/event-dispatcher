package response

import (
	"encoding/json"
	"log"
	"strconv"
)

type LivescoreData struct {
	EventID       string   `json:"Eid"`
	EventStatus   string   `json:"Eps"`
	Team1Info     []TeamInfo `json:"T1"`
	Team2Info     []TeamInfo `json:"T2"`
	Team1ScoreFT  string   `json:"Tr1OR"`
	Team2ScoreFT  string   `json:"Tr2OR"`
	EventFinished int      `json:"epr"` //epr=0 game not started, epr=1 game ongoing, epr=2 game finished
}

type TeamInfo struct {
	TeamName string `json:"Nm"`
}

func (p LivescoreData) GetTeam1Name() string {
	return p.Team1Info[0].TeamName
}

func (p LivescoreData) GetTeam2Name() string {
	return p.Team2Info[0].TeamName
}

func (p LivescoreData) GetTeamHomeScore() int {
	n, _ := strconv.Atoi(p.Team1ScoreFT)
	return n
}

func (p LivescoreData) GetTeamAwayScore() int {
	n, _ := strconv.Atoi(p.Team2ScoreFT)
	return n
}

func (p LivescoreData) Unmarshal(bytes []byte) any {
	var output LivescoreData
	if err := json.Unmarshal(bytes, &output); err != nil {
		log.Fatal(err)
	}
	return output
}