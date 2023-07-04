package requestSource

type Payload interface {
	LivescoreData | DummyData
}

type RequestAccessor[T Payload] interface {
	GetData() T
	GetTeam1Name() string
	GetTeam2Name() string
	GetTeam1Score() int
	GetTeam2Score() int
}