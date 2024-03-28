package poller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kamilszymczak/event-dispatcher/response"
)

func TestPoolingDataFromLivescore(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"Eid":"909663","Tr1OR":"2","Tr2OR":"4","T1":[{"Nm":"Pomigliano Women"}],"T2":[{"Nm":"Sampdoria Women"}],"Eps":"FT","Epr":2}`))
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
