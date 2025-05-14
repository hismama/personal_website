package main

import (
	"fmt"
	"github.com/hismama/personal_website/handlers"
	"github.com/hismama/personal_website/handlers/game"
	"github.com/hismama/personal_website/templates"
	"log"
	"net/http"
	"os"
)

func init() {
	game.LoadGlobalStats()
	templates.Init()
}

func main() {
	// Pages
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)
	http.HandleFunc("/resume", handlers.Resume)
	http.HandleFunc("/projects", handlers.Projects)
	http.HandleFunc("/game", handlers.Game)
	http.HandleFunc("/contact", handlers.Contact)

	// Assets
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Game
	http.HandleFunc("/click", game.ClickHandler)
	http.HandleFunc("/start-race", game.StartRaceHandler)
	http.HandleFunc("/end-race", game.EndRaceHandler)
	http.HandleFunc("/score", game.ScoreBoard)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
