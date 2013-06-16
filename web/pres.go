package web

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

func NewPresHandle(path string) (*PresHandle, error) {
	raw, err := ioutil.ReadFile(path)

	result := &PresHandle{
		presentation: string(raw),
	}
	return result, err
}

var code = "asd123"

type PresHandle struct {
	presentation string
}

func (pres *PresHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		pres.getPresHandle(w, r)
		return
	}

	if r.Method == "POST" {
		pres.postPresHandle(w, r)
		return
	}

	http.NotFound(w, r)
}

func (pres *PresHandle) postPresHandle(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, cookieName)

	r.ParseForm()

	c := r.FormValue("code")

	if c == code {
		session.Options.MaxAge = 0
		store.Save(w, r, session)
	}

	http.Redirect(w, r, "/", http.StatusFound)

}

func (pres *PresHandle) getPresHandle(w http.ResponseWriter, r *http.Request) {
	_, isNew := store.Get(r, cookieName)

	w.Header().Add("content-type", "text/html")

	if isNew {
		templates.ExecuteTemplate(w, "login", nil)
		return
	}
	values := make(map[string]interface{})

	values["pres"] = pres.presentation

	name, err := os.Hostname()

	if err != nil {
		log.Println(err)
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		log.Println(err)
	}

	values["url"] = fmt.Sprint("ws://", addrs[0], ":8080/ws")

	templates.ExecuteTemplate(w, "index", values)
	return
}
