package response

import "github.com/kamilszymczak/event-dispatcher/request"

// Type Constraint
type Payload interface {
	LivescoreData | DummyData
	ResponseAccessor
}

type Response[T Payload] struct {
	request  request.Request
	response T
}

type Responser interface {
	ResponseAccessor
}

type ResponseAccessor interface {
	GetTeam1Name() string
	GetTeam2Name() string
	GetTeamHomeScore() int
	GetTeamAwayScore() int
}

func New[T Payload](req request.Request, res T) *Response[T] {
	return &Response[T]{
		request:  req,
		response: res,
	}
}

func (r *Response[T]) GetData() ResponseAccessor {
	return r.response
}
