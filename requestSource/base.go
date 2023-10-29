package requestSource

// Type Constraint
type Payload interface {
	LivescoreData | DummyData
}

type Sourceable[T Payload] interface {
	GetData() T
	SetData(T) T
}

type ResponseAccessor interface {
	GetTeam1Name() string
	GetTeam2Name() string
	GetTeam1Score() int
	GetTeam2Score() int
}
