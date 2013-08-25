package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"fmt"
	"fp-shell/parse"
	"fp-shell/socket"
	"fp-shell/web"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

type Presentation struct {
	Title       string
	CustomStyle []string
	Slides      []*Slide
}

type Slide struct {
	Title   string
	Class   []string
	Sclass  []string
	Content string
	Notes   string
	Ain     string
	Aout    string
}

func main() {
	flag.Parse()

	var pres *web.PresHandle = nil

	if len(os.Args) == 1 {
		log.Println("Kein File spezifiziert")
		return
	}

	if strings.HasSuffix(os.Args[1], ".fp") {

		raw, err := ioutil.ReadFile(os.Args[1])

		if err != nil {
			panic(err)
			return
		}

		json := parse.NewJSONOutput(raw)

		pres = web.NewJSONPresHandle(json)

	} else {

		var err error
		pres, err = web.NewPresHandle(os.Args[1])

		if err != nil {
			log.Println(err)
		}
	}

	http.Handle("/", pres)

	ctrl, err := web.NewControllerHandle(os.Args[1])

	if err != nil {
		log.Println(err)
	}
	http.Handle("/ctrl", ctrl)

	rootpath := os.Getenv("FPROOT")

	http.HandleFunc("/j/", web.StaticFileHandler(fmt.Sprint(rootpath, "/static/")))
	http.HandleFunc("/s/", web.StaticFileHandler(fmt.Sprint(rootpath, "/static/")))

	server := socket.NewServer()

	http.Handle("/ws", websocket.Handler(server.WebsocketHandler))

	displayInfo()

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}

}

func displayInfo() {
	fmt.Printf("Code: %s\n", "asd123")

	name, err := os.Hostname()

	if err != nil {
		log.Println(err)
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("IP-Host: %s\n", addrs[0])
}
