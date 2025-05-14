package game

import (
	"cloud.google.com/go/firestore"
	"encoding/json"
	"github.com/hismama/personal_website/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"sync"
	"time"
)

type GlobalStatsStruct struct {
	TotalClicks  int `firestore:"total_clicks"`
	SiteHighRace int `firestore:"site_high_race"`
}

var GlobalStats GlobalStatsStruct

var (
	clickBuffer    int
	buffMutex      sync.Mutex
	flushThreshold = 20
	flushInterval  = 30 * time.Second
	docRef         = config.FS.Doc("stats/global")
)

func init() {
	go startClickFlusher()
}

func startClickFlusher() {
	ticker := time.NewTicker(flushInterval)
	for _ = range ticker.C {
		flushClicks()
	}
}

func ScoreBoard(w http.ResponseWriter, r *http.Request) {
	s, _ := Store.Get(r, "gopher-session")
	json.NewEncoder(w).Encode(map[string]int{
		"score":       s.Values["score"].(int),
		"sessionHigh": s.Values["race_high"].(int),
		"global_high": GlobalStats.SiteHighRace,
		"total":       GlobalStats.TotalClicks,
	})
}

func LoadGlobalStats() {
	snap, err := docRef.Get(config.Ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			_, _ = docRef.Set(config.Ctx, GlobalStatsStruct{})
		}
		return
	}
	_ = snap.DataTo(&GlobalStats)
}

func incrementClicks(n int) {
	buffMutex.Lock()
	defer buffMutex.Unlock()
	clickBuffer += n
	if clickBuffer > flushThreshold {
		go flushClicks()
	}
}

func flushClicks() {
	buffMutex.Lock()
	n := clickBuffer
	clickBuffer = 0
	buffMutex.Unlock()

	if n == 0 {
		return
	}

	GlobalStats.TotalClicks += n

	_, err := docRef.Update(config.Ctx, []firestore.Update{
		{Path: "total_clicks", Value: firestore.Increment(n)},
	})
	if status.Code(err) == codes.NotFound {
		_, _ = docRef.Set(config.Ctx, GlobalStatsStruct{
			TotalClicks:  GlobalStats.TotalClicks,
			SiteHighRace: GlobalStats.SiteHighRace,
		})
	} else if err != nil {
		log.Printf("incrementClicks: %v", err)
	}
}

func saveGlobalHigh(n int) {
	GlobalStats.SiteHighRace = n
	_, err := docRef.Update(config.Ctx, []firestore.Update{
		{Path: "site_high_race", Value: GlobalStats.SiteHighRace},
	})
	if status.Code(err) == codes.NotFound {
		_, _ = docRef.Set(config.Ctx, GlobalStatsStruct{
			TotalClicks:  GlobalStats.TotalClicks,
			SiteHighRace: GlobalStats.SiteHighRace,
		})
	} else if err != nil {
		log.Printf("saveGlobalHigh: %v", err)
	}
}
