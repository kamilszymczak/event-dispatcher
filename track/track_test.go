package track

import (
	"github.com/kamilszymczak/event-dispatcher/request"
	"github.com/kamilszymczak/event-dispatcher/requestSource"
	"testing"
	"time"
)

func TestTrackRequestNoRefreshRate(t *testing.T) {
	req := request.New[requestSource.LivescoreData]("https://prod-public-api.livescore.com/v1/api/app/scoreboard/soccer/909663")

	tracker, _ := New[requestSource.LivescoreData]()
	tracker.RefreshRate(4 * time.Minute)
	tracker.AddRequest(req)

	if tracker.requests[0].GetRefreshRate() != tracker.refreshRate {
		t.Errorf("Expected refresh rate of %s, got %s", tracker.refreshRate, tracker.requests[0].GetRefreshRate())
	}
}

func TestTrackListenOnce(t *testing.T) {

}

func TestTrackComputeFunc(t *testing.T) {

}
