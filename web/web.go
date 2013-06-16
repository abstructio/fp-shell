package web

import (
	"fmt"
	"github.com/ioboi/sessions"
	"html/template"
	"net/http"
)

const cookieName = "pressession"

var store = sessions.NewMemorySessionStore()

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob("template/*"))
}

func StaticFileHandler(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprint(path, r.URL.Path[3:]))
	}
}
