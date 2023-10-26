package requestSource

// type LivescoreAccessor interface {
// 	RequestAccessor[LivescoreData]
// 	GetTeam1Name() string
// 	GetTeam2Name() string
// 	GetTeam1Score() int
// 	GetTeam2Score() int
// }

type LivescoreData struct {
	EventID string `json:"Eid"`
	EventStatus string `json:"Eps"`
	Team1ScoreFT string `json:"Tr1OR"`
	Team2ScoreFT string `json:"Tr2OR"`
	EventFinished int `json:"epr"` //epr=0 game not started, epr=1 game ongoing, epr=2 game finished
}

func (p *LivescoreData) GetData() *LivescoreData {
	return p
}