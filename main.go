package main

import (
	"fmt"

	"github.com/kamilszymczak/event-dispatcher/request"
	"github.com/kamilszymczak/event-dispatcher/requestSource"
	"github.com/kamilszymczak/event-dispatcher/track"
)

func main() {
	// refresh rate should restart rate countdown again once request is received, not strictly every x interval
	// since getting read from api could take time, and we don't want queued 'messages' to the api

	// Produce - fetch from api
	req1 := request.New[requestSource.LivescoreData]("https://prod-public-api.livescore.com/v1/api/app/scoreboard/soccer/909663")
	req2 := request.New[requestSource.LivescoreData]("https://prod-public-api.livescore.com/v1/api/app/scoreboard/soccer/714150")

	tracker, _ := track.New[requestSource.LivescoreData]()
	tracker.AddRequest(req1, req2)

	// Returns a read only channel (<- chan T). Possibly use observer pattern
	//tracker.SetComputeFunc(computeScoreChanged)
	trackCh := tracker.Listen()

	// Publish
P:
	for {
		select {
		case req, ok := <-trackCh:
			if !ok {
				fmt.Println("Publish channel closing")
				break P
			}
			fmt.Println("received: ", req.GetData())
		}
	}

	fmt.Println("Main: Ending")
}

// use type that can be compared and compare
func computeScoreChanged(current, new int) bool {
	return current != new
}
