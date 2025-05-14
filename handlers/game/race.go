package game

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func StartRaceHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "gopher-session")
	session.Values["race_start"] = time.Now()
	session.Values["race_score"] = 0
	session.Values["spawn_page"] = "/game"
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func EndRaceHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "gopher-session")
	raceScore := session.Values["race_score"].(int)

	if raceHigh, ok := session.Values["race_high"].(int); !ok || raceScore > raceHigh {
		session.Values["race_high"] = raceScore
	}
	if raceScore > GlobalStats.SiteHighRace {
		saveGlobalHigh(raceScore)
	}
	delete(session.Values, "race_start")
	delete(session.Values, "race_score")
	session.Save(r, w)

	json.NewEncoder(w).Encode(map[string]int{
		"raceScore":    raceScore,
		"sessionHigh":  session.Values["race_high"].(int),
		"globalHigh":   GlobalStats.SiteHighRace,
		"totalCatches": GlobalStats.TotalClicks,
	})
}
