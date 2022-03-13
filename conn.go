package ynot

import (
	"fmt"
	"io"
	"log"
	"net"
)

type Conn interface {
	// Start activate the connection.
	Start()

	// Stop close the connection.
	Stop()

	// GetTCPConn gets the underlying TCP connection.
	GetTCPConn() *net.TCPConn

	// GetConnID gets the unique ID of TCP connection.
	GetConnID() uint32

	// RemoteAddr returns the remote address.
	RemoteAddr() net.Addr

	// Send sends data to the client.
	Send(data []byte) error
}

// The HandlerFunc represent the business logic which connection banding.
type HandlerFunc func(*net.TCPConn, []byte, int) error

type YConn struct {
	// TCP connection.
	Conn *net.TCPConn

	// The unique ID of TCP connection.
	ConnID uint32

	// isClosed represents whether the connection is closed.
	isClosed bool

	// the banding business logic of connection.
	handlerFunc HandlerFunc

	// Channel for notifying close.
	StopChan chan struct{}
}

func NewConn(conn *net.TCPConn, connID uint32, handlerFunc HandlerFunc) Conn {
	return &YConn{
		Conn:        conn,
		ConnID:      connID,
		handlerFunc: handlerFunc,
		StopChan:    make(chan struct{}),
	}
}

func DefaultHandlerFunc(conn *net.TCPConn, data []byte, n int) error {
	if _, err := conn.Write(data[:n]); err != nil {
		return fmt.Errorf("DefaultHandler write error: %s", err)
	}
	return nil
}

func (yc *YConn) Start() {
	log.Println("Conn start, ConnID =", yc.ConnID)

	// Start a goroutine to handle reading business
	go yc.readerPipeline()

	// Start a goroutine to handle writing business

}

func (yc *YConn) readerPipeline() {
	log.Println("ConnID =", yc.ConnID, ", reader goroutine is running")

	defer func() {
		log.Println("ConnID =", yc.ConnID, ", reader goroutine exits, remote addr is", yc.RemoteAddr().String())
		yc.Stop()
	}()

	for {
		buf := make([]byte, 512)
		n, err := yc.Conn.Read(buf)
		if err == io.EOF {
			log.Println("Client close conn, remote addr is", yc.RemoteAddr().String())
			break
		}
		if err != nil {
			log.Println("ConnID =", yc.ConnID, ", read buf error:", err)
			continue
		}

		// Invoke corresponding handler function.
		err = yc.handlerFunc(yc.Conn, buf, n)
		if err != nil {
			log.Println("ConnID =", yc.ConnID, ", handler func error:", err)
			continue
		}
	}
}

func (yc *YConn) Stop() {
	log.Println("Conn stop... ConnID =", yc.ConnID)

	if yc.isClosed {
		return
	}

	// Close TCP connection.
	yc.Conn.Close()

	// Close the channel.
	close(yc.StopChan)

	// Set closing state.
	yc.isClosed = true
}

func (yc *YConn) GetTCPConn() *net.TCPConn {
	return yc.Conn
}

func (yc *YConn) GetConnID() uint32 {
	return yc.ConnID
}

func (yc *YConn) RemoteAddr() net.Addr {
	return yc.Conn.RemoteAddr()
}

func (yc *YConn) Send(data []byte) error {
	return nil
}
