package main

import (
	"fmt"
	"log"

	"github.com/kamilszymczak/event-dispatcher/poller"
	"github.com/kamilszymczak/event-dispatcher/response"
)

func main() {
	// refresh rate should restart rate countdown again once request is received, not strictly every x interval
	// since getting read from api could take time, and we don't want queued 'messages' to the api

	// Produce - fetch from api
	pollerService := poller.New("https://prod-public-api.livescore.com/v1/api/app/scoreboard/soccer/", response.LivescoreData{})
	pollerService.AddObservable(poller.Observable{Address: "909663"})

	listenChan := pollerService.Listen()

	P:
	for {
		select {
		case req, ok := <-listenChan:
			if !ok {
				fmt.Println("Publish channel closing")
				break P
			}
			typedRes, _ := pollerService.HandleEvent(req)
			log.Print("received: ", typedRes)
		}
	}

	fmt.Println("Main: Ending")
}
