package game

import (
	"encoding/gob"
	"github.com/gorilla/sessions"
	"github.com/hismama/personal_website/config"
	"math/rand"
	"time"
)

func init() {
	gob.Register(time.Time{})
}

var Store = sessions.NewCookieStore(config.CookieSecret)

var SitePages = []string{"/", "/about", "/resume", "/projects", "/contact", "/game"}

// Timer Circumference
const full = 300.0

type SpawnInfo struct {
	Spawn       bool
	RaceActive  bool
	RaceScore   int
	TimeLeft    int
	TimerOffset int
	TimerColor  string
	Top, Left   int
}

func GopherInfo(s *sessions.Session, path string) SpawnInfo {
	raceActive, remain, raceScore := false, 0, 0
	if raw, ok := s.Values["race_start"]; ok {
		elapsed := time.Since(raw.(time.Time))
		if elapsed < 30*time.Second {
			raceActive = true
			remain = int((30*time.Second - elapsed).Seconds())
		}
	}

	if v := s.Values["race_score"]; v != nil {
		raceScore = v.(int)
	}

	if s.IsNew {
		next := SitePages[rand.Intn(len(SitePages))]
		s.Values["spawn_page"] = next
	}

	target, _ := s.Values["spawn_page"].(string)
	shouldSpawn := target == path
	//Colors line up from CSS on circle
	var color string
	switch {
	case raceActive && remain <= 5:
		color = "#dc2626"
	case raceActive && remain <= 15:
		color = "#f97316"
	default:
		color = "#22c55e"
	}

	return SpawnInfo{
		Spawn:       shouldSpawn,
		RaceActive:  raceActive,
		RaceScore:   raceScore,
		TimeLeft:    remain,
		TimerOffset: int(full * (1 - float64(remain)/30)),
		TimerColor:  color,
		Top:         rand.Intn(71) + 10,
		Left:        rand.Intn(71) + 10,
	}
}
