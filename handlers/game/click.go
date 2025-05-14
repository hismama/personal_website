package game

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type clickReq struct {
	Path string `json:"path"`
}

type clickResp struct {
	Score        int    `json:"score"`
	TotalCatches int    `json:"totalCatches"`
	RaceScore    int    `json:"raceScore"`
	Message      string `json:"message"`
	Gopher       string `json:"gopher"`
	FunFact      string `json:"funFact"`
	SpawnAgain   bool   `json:"spawnAgain"`
}

var funFacts = []string{
	"2019‑2023: turned QA automation chops into full‑stack dev skills.",
	"Built a concurrent EPCIS 1.2 XML generator in Go (>100 K serials/sec).",
	"Scaled a Django REST API with selective update_fields saves.",
	"Front‑end minimalism fan: pure HTML/CSS + Go’s net/http templates.",
	"Runs “Misheard Tees”: print‑on‑demand store for lyric puns.",
	"Loves refactoring: fewer duplicate selectors & smarter helpers.",
	"Go mantra: “make it obvious, then make it blazing fast.”",
	"Built 150+ modular XML validations with toggles, templates, and severity levels.",
	"Wrote a serial import engine to avoid double-claiming during platform migrations.",
	"Technically a rocket scientist. Still debugged CSS like everyone else.",
	"Learned orbital mechanics. Now everything revolves around JSON.",
	"Holds a B.S. in Aerospace Engineering with a strong foundation in systems thinking.",
	"Licensed SCUBA diver—comfortable under pressure, literally.",
	"Used SolidWorks and TIG welding to bring ideas to life in the machine shop.",
	"Believes JavaScript is for behavior, not for putting boxes where CSS belongs.",
	"Advocates mastering core web fundamentals (HTML/CSS) before diving into frameworks.",
	"Built this site so the layout stays intact with or without JavaScript—only the fun stuff breaks.",
}

func ClickHandler(w http.ResponseWriter, r *http.Request) {
	var body clickReq
	_ = json.NewDecoder(r.Body).Decode(&body)
	currentPath := body.Path
	session, _ := Store.Get(r, "gopher-session")

	if session.Values["score"] == nil {
		session.Values["score"] = 0
	}
	if session.Values["race_score"] == nil {
		session.Values["race_score"] = 0
	}

	points, message, gopherType := 1, "Caught a regular gopher!", "normal"

	// Normal score always increments
	session.Values["score"] = session.Values["score"].(int) + points

	// If race active (<30s window) add to race_score too
	if startRaw, ok := session.Values["race_start"]; ok && time.Since(startRaw.(time.Time)) < 30*time.Second {
		session.Values["race_score"] = session.Values["race_score"].(int) + points
	}

	spawnPage := SitePages[rand.Intn(len(SitePages))]
	session.Values["spawn_page"] = spawnPage
	spawnImmediate := spawnPage == currentPath

	session.Save(r, w)
	incrementClicks(points)

	fact := funFacts[rand.Intn(len(funFacts))]

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(clickResp{
		Score:        session.Values["score"].(int),
		TotalCatches: GlobalStats.TotalClicks,
		RaceScore:    session.Values["race_score"].(int),
		Message:      message,
		Gopher:       gopherType,
		FunFact:      fact,
		SpawnAgain:   spawnImmediate,
	})
	if err != nil {
		return
	}
}
