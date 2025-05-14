package handlers

import (
	"github.com/hismama/personal_website/utils"
	"net/http"
)

func Resume(w http.ResponseWriter, r *http.Request) {
	err := utils.RenderTemplate(w, r, "resume.gohtml", nil)
	utils.Check(err, w)
}
