package main

import (
	"fmt"

	"github.com/kamilszymczak/event-dispatcher/request"
	"github.com/kamilszymczak/event-dispatcher/requestSource"
	"github.com/kamilszymczak/event-dispatcher/track"
)

func main() {	
	// refreshrate should restart rate countdown again once request is received, not strictly every x interval
	// since getting read from api could take time and we don't want queued 'messages' to the api

	// Produce - fetch from api
	var requests []request.Requestable[requestSource.LivescoreData]

	req1 := request.New[requestSource.LivescoreData]("https://prod-public-api.livescore.com/v1/api/app/scoreboard/soccer/909663")
	req2 := request.New[requestSource.LivescoreData]("https://prod-public-api.livescore.com/v1/api/app/scoreboard/soccer/714150")
	requests = append(requests, req1, req2)


	tracker, _ := track.New[requestSource.LivescoreData]()
	tracker.AddRequest(req1, req2)

	// Returns a read only channel (<- chan T). Possibly use observer pattern
	trackCh := tracker.Listen()

	// Publish
	P: for {
		select {
		case req, ok := <-trackCh:
			if(!ok){
				break P
			}
			fmt.Println("received: ", req)

		}
	}

	// time.Sleep(6 * time.Second)
    // fmt.Println("Goodye!! to Main function")
}