package handlers

import (
	"github.com/hismama/personal_website/utils"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	err := utils.RenderTemplate(w, r, "index.gohtml", nil)
	utils.Check(err, w)
}
