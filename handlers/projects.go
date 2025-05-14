package handlers

import (
	"github.com/hismama/personal_website/utils"
	"net/http"
)

func Projects(w http.ResponseWriter, r *http.Request) {
	err := utils.RenderTemplate(w, r, "projects.gohtml", nil)
	utils.Check(err, w)
}
