package ynot

import (
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestYServer(t *testing.T) {
	ip := "127.0.0.1"
	port := 8080
	network := "tcp"
	num := 100
	// Start a server.
	s := New("ynot-server", ip, network, port)
	go func() {
		s.Serve()
	}()
	ss := s.(*YServer)
	// Wait for server running.
	for !ss.IsRunning() {
		time.Sleep(time.Microsecond)
	}
	// Start a client to test the functionality of the server.
	conn, err := net.Dial(network, fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		t.Fatal(err)
	}
	base := "hello"
	for i := 0; i < num; i++ {
		body := base + strconv.Itoa(i)
		n, err := conn.Write([]byte(body))
		if err != nil {
			t.Error(err)
		}
		// Compare sent content and received content
		b := make([]byte, 512)
		n, err = conn.Read(b)
		if err != nil {
			t.Error(err)
		}
		if string(b[:n]) != body {
			t.Errorf("compare error, expected: %s， got: %s", body, string(b[:n]))
		}
	}
}
