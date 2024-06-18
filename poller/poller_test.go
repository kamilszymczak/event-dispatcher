package poller

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/kamilszymczak/event-dispatcher/response"
)

func TestPoolingDataFromLivescore(t *testing.T) {
	rawJson, _ := os.ReadFile("./poller_livescoreOutputValid_test.json")
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(rawJson))
	}))
	defer server.Close()

	expected := response.LivescoreData{
		EventID: "909663",
		EventStatus: "FT",
		Team1Info: []response.TeamInfo{
			{
				TeamName: "Pomigliano Women",
			},
		},
	}

	poller := &poller{apiUrl: server.URL, responseType: response.LivescoreData{}}

	obs := Observable{Address: ""}

	data, _ := poller.fetchData(obs)
	event := buildEvent(obs, data)
	got, err := poller.HandleEvent(event)

	if err != nil {
		t.Errorf("event handled with error: %t", err)
	}

	if got.GetTeam1Name() != expected.GetTeam1Name() {
		t.Errorf("Expected %s got %s", expected.GetTeam1Name(), got.GetTeam1Name())
	}
}

func TestPoolingDataFromLivescoreIncorrectJSONFormat(t *testing.T) {
	// Service returns invalid JSON format, missing Tr1OR
	rawJson, _ := os.ReadFile("./poller_livescoreOutputInvalid_test.json")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(rawJson))
	}))
	defer server.Close()

	expected := response.LivescoreData{
		EventID: "909663",
		EventStatus: "FT",
		Team1Info: []response.TeamInfo{
			{
				TeamName: "",
			},
		},
		Team1ScoreFT: "0",
	}

	poller := &poller{apiUrl: server.URL, responseType: response.LivescoreData{}}

	obs := Observable{Address: ""}

	data, _ := poller.fetchData(obs)
	event := buildEvent(obs, data)
	got, err := poller.HandleEvent(event)

	if err != nil {
		t.Errorf("event handled with error: %t", err)
	}

	if got.GetTeam1Name() != expected.GetTeam1Name() {
		t.Errorf("Expected %s got %s", expected.GetTeam1Name(), got.GetTeam1Name())
	}

	if got.GetTeamHomeScore() != expected.GetTeamHomeScore() {
		t.Errorf("Expected %v got %v", expected.GetTeamHomeScore(), got.GetTeamHomeScore())
	}
}

func TestPoolingDataFromLivescoreInvalidAddress(t *testing.T) {
	poller := &poller{apiUrl: "https://prod-public-api.livescore.com/v1/api/app/scoreboard/soccer/", responseType: response.LivescoreData{}}
	obs := Observable{Address: "108583400"}
	want := 410

	_, err := poller.fetchData(obs)

	if !strings.Contains(err.Error(), strconv.Itoa(want)) {
		t.Errorf("Expected error to contain status code %v got %v", 410, err.Error())
	}
}

func TestDispatchFunc(t *testing.T) {
	dispatchFunc := func(resp response.ResponseAccessor) bool {
		return resp.GetGameStatus() == 2
	}

	rawJson, _ := os.ReadFile("./poller_livescoreOutputValid_test.json")
	poller := &poller{apiUrl: "", responseType: response.LivescoreData{}}
	poller.SetDispatchFunc(dispatchFunc)

	event := buildEvent(Observable{}, rawJson)
	response, _ := poller.HandleEvent(event)
	isMatchOver := poller.dispatchFunc(response)

	if isMatchOver != true {
		t.Errorf("Expected %v got %v", true, isMatchOver)
	}
}

func TestDontDispatch(t *testing.T) {
	dispatchFunc := func(resp response.ResponseAccessor) bool {
		return resp.GetGameStatus() == 1
	}

	rawJson, _ := os.ReadFile("./poller_livescoreOutputValid_test.json")
	poller := &poller{apiUrl: "", responseType: response.LivescoreData{}}
	poller.SetDispatchFunc(dispatchFunc)

	event := buildEvent(Observable{}, rawJson)
	response, _ := poller.HandleEvent(event)
	matchOngoing := poller.dispatchFunc(response)

	if matchOngoing != false {
		t.Errorf("Expected %v got %v", false, matchOngoing)
	}
}