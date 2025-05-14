package templates

import (
	"github.com/dustin/go-humanize"
	"html/template"
)

var Tpl *template.Template

var funcMap = template.FuncMap{
	"comma": CommaInt,
}

func Init() {
	Tpl = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/helpers/*.gohtml"))
	Tpl = template.Must(Tpl.ParseGlob("templates/base/*.gohtml"))
}

func CommaInt(num int) string {
	return humanize.Comma(int64(num))
}
