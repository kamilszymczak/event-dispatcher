package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/kamilszymczak/event-dispatcher/request"
	"github.com/kamilszymczak/event-dispatcher/requestSource"
)

func main() {	
	// Run tracker with .listen() command which returns a read only channel (<- chan T)
	// and use observer pattern, .listen() is generator pattern.
	// <- tracker.listen()

	// refreshrate should restart rate countdown again once request is received, not striclty every x interval
	// since getting read from api could take time and we don't want queued 'messages'

	// get data from listen tracker to observer and use decoration pattern to check if same

	// Produce - fetch from api
	var requests []request.Requestable[requestSource.LivescoreData]

	req1 := request.New[requestSource.LivescoreData]("https://prod-public-api.livescore.com/v1/api/app/scoreboard/soccer/909663")
	req2 := request.New[requestSource.LivescoreData]("https://prod-public-api.livescore.com/v1/api/app/scoreboard/soccer/714150")
	requests = append(requests, req1, req2)
	listenerCh := trackerToChannel[requestSource.LivescoreData](requests)

	// tracker := tracker.New()
	// tracker.AddRequest(req1)
	// tracker.AddRequest(req2)
	// trackCh <- tracker.listen()

	// Publish
	P: for {
		select {
		case req, ok := <-listenerCh:
			if(!ok){
				break P
			}
			fmt.Println("received: ", req)

		}
	}

	time.Sleep(6 * time.Second)
    fmt.Println("Goodye!! to Main function")

	// Consume/Map - Read in json and check if different from previous

	// Publish - 
	// for _, out := range requests {
	// 	fmt.Println(out.GetData().EventID)
	// }
}

func fetchData[T requestSource.Payload](channel chan<- request.Requestable[T], i int, request request.Requestable[T]) {
	res, err := http.Get(request.GetUrl())

	if i < 1 {
		fmt.Println("sleeping 2 sec " , i)
		time.Sleep(2 * time.Second)
		fmt.Println("done sleeping")
	}


	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	var output T
	json.Unmarshal(body, &output)

	request.SetData(output)
	channel <- request
}

func trackerToChannel[T requestSource.Payload](requests []request.Requestable[T]) <-chan request.Requestable[T] {
	out := make(chan request.Requestable[T])

	for i, request := range requests {
		go fetchData(out, i, request)
	}

	return out
}

// func parseRequest(<-chan request.Requestable[requestSource.LivescoreData]) <-chan request.Requestable[requestSource.LivescoreData] {


// }