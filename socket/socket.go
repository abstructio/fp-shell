package socket

import (
	ws "code.google.com/p/go.net/websocket"
	"log"
)

type command struct {
	Id    int
	Value interface{}
}

type Server struct {
	machines   map[*Machine]bool
	broadcast  chan command
	register   chan *Machine
	unregister chan *Machine
}

func NewServer() *Server {
	server := &Server{
		machines:   make(map[*Machine]bool),
		broadcast:  make(chan command),
		register:   make(chan *Machine),
		unregister: make(chan *Machine),
	}

	go server.run()

	return server
}

func (s *Server) run() {
	for {
		select {
		case cmd := <-s.broadcast:
			for m, _ := range s.machines {
				m.send <- cmd
			}
		case m := <-s.unregister:
			close(m.send)
			delete(s.machines, m)
		case m := <-s.register:
			s.machines[m] = true
		}
	}
}

func (s *Server) WebsocketHandler(conn *ws.Conn) {
	m := &Machine{
		conn:   conn,
		send:   make(chan command),
		server: s,
	}

	s.register <- m
	go m.SendConnection()
	m.ReadConnection()
}

type Machine struct {
	conn   *ws.Conn
	send   chan command
	server *Server
}

func (m *Machine) SendConnection() {
	for cmd := range m.send {
		err := ws.JSON.Send(m.conn, cmd)
		if err != nil {
			log.Println(err)
		}
	}
}

func (m *Machine) ReadConnection() {
	for {
		var cmd command
		err := ws.JSON.Receive(m.conn, &cmd)
		if err != nil {
			log.Println(err)
			break
		}
		m.server.broadcast <- cmd
	}
	m.server.unregister <- m
}
