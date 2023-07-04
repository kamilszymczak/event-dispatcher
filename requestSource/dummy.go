package requestSource

type DummyAccessor interface {
	RequestAccessor[DummyData]
}

type DummyData struct {
	EventID string `json:"Eid"`
	EventStatus string `json:"Eps"`
	Team1ScoreFT string `json:"Tr1OR"`
	Team2ScoreFT string `json:"Tr2OR"`
	EventFinished string `json:"epr"` //epr=0 game not started, epr=1 game ongoing, epr=2 game finished
	EventDummy bool `json:"dummy"`
}

func (p *DummyData) GetData() *DummyData {
	return p
}