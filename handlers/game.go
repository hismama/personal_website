package handlers

import (
	"github.com/hismama/personal_website/handlers/game"
	"github.com/hismama/personal_website/utils"
	"net/http"
)

func Game(w http.ResponseWriter, r *http.Request) {
	session, _ := game.Store.Get(r, "gopher-session")
	raceHigh := 0
	if v := session.Values["race_high"]; v != nil {
		raceHigh = v.(int)
	}
	data := struct {
		RaceHigh int
		Global   int
	}{
		RaceHigh: raceHigh,
		Global:   game.GlobalStats.SiteHighRace,
	}

	err := utils.RenderTemplate(w, r, "game.gohtml", data)
	utils.Check(err, w)
}
