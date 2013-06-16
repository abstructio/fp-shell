package main

import (
	"flag"
	"fmt"
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

	http.HandleFunc("/j/", web.StaticFileHandler("static/"))
	http.HandleFunc("/s/", web.StaticFileHandler("static/"))

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
