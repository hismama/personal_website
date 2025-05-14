package handlers

import (
	"github.com/hismama/personal_website/utils"
	"net/http"
)

func About(w http.ResponseWriter, r *http.Request) {
	err := utils.RenderTemplate(w, r, "about.gohtml", nil)
	utils.Check(err, w)
}
