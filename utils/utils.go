package utils

import (
	"github.com/hismama/personal_website/handlers/game"
	"github.com/hismama/personal_website/templates"
	"log"
	"net/http"
	"os"
)

var prod = os.Getenv("ENV") == "prod"

type TemplateData struct {
	Prod         bool
	Path         string
	Score        int
	TotalCatches int
	Spawn        game.SpawnInfo
	Data         interface{}
}

func Check(err error, w http.ResponseWriter) {
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) error {
	tpl, err := templates.Tpl.Clone()
	if err != nil {
		return err
	}
	tpl, err = tpl.ParseFiles(
		"templates/pages/" + tmpl,
	)
	if err != nil {
		return err
	}
	path := r.URL.Path
	session, _ := game.Store.Get(r, "gopher-session")
	score := 0
	if v := session.Values["score"]; v != nil {
		score = v.(int)
	}

	var spawn = game.SpawnInfo{}
	spawn = game.GopherInfo(session, path)
	_ = session.Save(r, w)

	td := TemplateData{
		Prod:         prod,
		Path:         path,
		Score:        score,
		TotalCatches: game.GlobalStats.TotalClicks,
		Spawn:        spawn,
		Data:         data,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tpl.ExecuteTemplate(w, tmpl, td)
	if err != nil {
		return err
	}

	return nil
}
