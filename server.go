package ynot

import (
	"fmt"
	"log"
	"net"
)

type serverState int

const (
	StateInit serverState = iota
	StateStart
	StateRunning
	StateCLosed
)

type Server interface {
	// Start starts tcp server.
	Start()
	// Serve implements business logic.
	Serve()
	// Stop stops tcp server.
	Stop()
}

// YServer is a struct which implement interface Server.
type YServer struct {
	// A server name.
	Name string
	// The IP that service listens on.
	IP string
	// The port that service listens on.
	Port int
	// Network type, such as tcp, tcp4, tcp6.
	Network string
	// A server state.
	state serverState
}

func (ys *YServer) Start() {
	log.Printf("Server start, listen on %s:%d\n", ys.IP, ys.Port)
	go func() {
		addr := fmt.Sprintf("%s:%d", ys.IP, ys.Port)
		ln, err := net.Listen(ys.Network, addr)
		if err != nil {
			log.Panicf("Server[%s:%d] start error: %s", ys.IP, ys.Port, err.Error())
		}
		ys.state = StateRunning
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Printf("Server[%s:%d] accept error: %s\n", ys.IP, ys.Port, err.Error())
				continue
			}

			yConn := NewConn(conn.(*net.TCPConn), 1, DefaultHandlerFunc)
			yConn.Start()
		}
	}()
}

func (ys *YServer) IsRunning() bool {
	return ys.state == StateRunning
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("Read buf error:", err)
			break
		}
		log.Println("Server receive:", string(buf[:n]))
		n, err = conn.Write(buf[:n])
		if err != nil {
			log.Println("Write buf error:", err)
		}
	}
}

func (ys *YServer) Serve() {
	// Start server.
	ys.Start()
	// TODO Complete business logic

	// Block server.
	select {}
}

func (ys *YServer) Stop() {
	//TODO implement me
	panic("implement me")
}

func New(name, ip, network string, port int) Server {
	return &YServer{
		Name:    name,
		IP:      ip,
		Network: network,
		Port:    port,
	}
}
