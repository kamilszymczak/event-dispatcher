package response

type Responser interface {
	ResponseAccessor
	Unmarshal(bytes []byte) any
}

type ResponseAccessor interface {
	GetTeam1Name() string
	GetTeam2Name() string
	GetTeamHomeScore() int
	GetTeamAwayScore() int
}

// func Insert[T Payload](req request.Request, res T) (*Response[T], bool) {
// 	journal := make(map[request.RequestableRefreshRater]*Response[T])

// 	_, ok := journal[&req]

// 	if !ok {
// 		newResponse := &Response[T]{
// 			request:  req,
// 			response: res,
// 		}
// 		journal[&req] = newResponse
// 		return journal[&req], true
// 	} else {
// 		// TODO: Compare if new differs from old if does then replace with new
// 		if responseDataChanged(journal[&req].response, res) {
// 			journal[&req].response = res
// 			return journal[&req], true
// 		}
// 	}
// 	return journal[&req], false
// }


// use type that can be compared and compare
// func responseDataChanged[T Payload](current, new T) bool {
// 	return current.GetTeamHomeScore() == new.GetTeamHomeScore() && current.GetTeamAwayScore() == new.GetTeamAwayScore()
// }
