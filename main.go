package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kamilszymczak/event-dispatcher/request"
	"github.com/kamilszymczak/event-dispatcher/entities"
)

func main() {
	// tracker, _ := track.New()
	// tracker.RefreshRate(time.Minute * 2)
	
	// request := request.NewRequest[ApiResult]("www.google.com", locateData)
	// tracker.AddRequest(request)

	// result := request.PerformLookup()
	// fmt.Println(result)

	// <- tracker.listen()

	// tracker.AddRequest[ApiResult]("www.google.com", locateScore)
	// tracker.AddRequest(&request)

	// Run tracker with .listen() command which returns a read only channel (<- chan T)
	// and use observer pattern, .listen() is generator pattern.

	// <- tracker.listen()

	// refreshrate should restart rate countdown again once request is received, not striclty every x interval
	// since getting read from api could take time and we don't want queued 'messages'

	// get data from listen tracker to observer and use decoration pattern to check if same


	// Produce - fetch from api
	var requests []request.Requestable[livescore.LivescoreData]

	req1 := request.New[livescore.LivescoreData]("https://prod-public-api.livescore.com/v1/api/app/scoreboard/soccer/909663")
	requests = append(requests, req1)

	// Consume/Map - Read in json and check if different from previous

	for _, req := range requests {
		// req := req.(request.Requestable[any])

		res, err := http.Get(req.GetUrl())
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
	
		body, err := ioutil.ReadAll(res.Body)
	
		var output livescore.LivescoreData
		json.Unmarshal(body, &output)

		req.SetData(output)
	}
	
	// Publish - 
	for _, out := range requests {
		fmt.Println(out.GetData().EventID)
	}

}