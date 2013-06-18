package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"fmt"
	"fp-shell/socket"
	"fp-shell/web"
	"log"
	"net"
	"net/http"
	"os"
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

	pres, err := web.NewPresHandle(os.Args[1])

	if err != nil {
		log.Println(err)
	}

	http.Handle("/", pres)

	ctrl, err := web.NewControllerHandle(os.Args[1])

	if err != nil {
		log.Println(err)
	}
	http.Handle("/ctrl", ctrl)

	http.HandleFunc("/j/", web.StaticFileHandler("static/"))
	http.HandleFunc("/s/", web.StaticFileHandler("static/"))

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
	    fmt.Printf("IP-Host: localhost\n")
		return
	}
	fmt.Printf("IP-Host: %s\n", addrs[0])
}
