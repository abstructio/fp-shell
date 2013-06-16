package web

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

type ControllerHandle struct {
	presentation string
}

func NewControllerHandle(path string) (*ControllerHandle, error) {
	raw, err := ioutil.ReadFile(path)
	result := &ControllerHandle{
		presentation: string(raw),
	}
	return result, err
}

func (ctrl *ControllerHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		ctrl.get(w, r)
		return
	}

	if r.Method == "POST" {
		ctrl.post(w, r)
		return
	}

	http.NotFound(w, r)
}

func (ctrl *ControllerHandle) get(w http.ResponseWriter, r *http.Request) {
	_, isNew := store.Get(r, cookieName)

	w.Header().Add("content-type", "text/html")

	if isNew {
		templates.ExecuteTemplate(w, "mobilecode", nil)
		return
	}

	values := make(map[string]interface{})

	values["pres"] = ctrl.presentation

	name, err := os.Hostname()

	if err != nil {
		log.Println(err)
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		log.Println(err)
	}

	values["url"] = fmt.Sprint("ws://", addrs[0], ":8080/ws")

	templates.ExecuteTemplate(w, "controller", values)
	return
}

func (ctrl *ControllerHandle) post(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, cookieName)

	r.ParseForm()

	c := r.FormValue("code")

	if c == code {
		session.Options.MaxAge = 0
		store.Save(w, r, session)
	}

	http.Redirect(w, r, "/ctrl", http.StatusFound)
}
