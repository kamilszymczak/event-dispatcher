package poller

import (
	"net/http"
	"net/http/httptest"
	"os"
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

	event := buildEvent(obs, poller.fetchData(obs))
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

	event := buildEvent(obs, poller.fetchData(obs))
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