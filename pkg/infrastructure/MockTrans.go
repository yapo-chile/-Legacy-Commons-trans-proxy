package infrastructure

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"
)

const (
	// WelcomeMessage is the message given when a client connects to the
	// server and can answer its request. It also signals the start of a
	// command response
	WelcomeMessage = "220 Welcome.\n"
	// BusyMessage is the message given when a client connects to the
	// server and the server is configured to refuse new connections, or
	// the number of pending connections is greater that the max connections
	BusyMessage = "521 Busy.\n"
	// EndMessage is message that define end of a command, or a
	// command response.
	EndMessage = "end\n"
)

// Handler function used to manage the received command and returns a response in bytes.
// The response must have the form of <param1>:<value1>\n<param2>:<value2>\n...,
// without the "end\n" message
type Handler func([]byte) []byte

// MockTransServer the struct tht represents a Mock trans server
type MockTransServer struct {
	Address  string
	IsBusy   bool
	listener net.Listener
	handler  Handler
	mtx      sync.RWMutex
}

// Start starts mock trans server
func (srv *MockTransServer) Start() {
	if srv.Address != "" {
		panic("trans test server already started")
	}
	srv.Address = srv.listener.Addr().String()
	go func() {
		_ = srv.Serve(srv.listener) // nolint: gosec
	}()
}

// Close shuts down the server and blocks until all outstanding
// requests on this server have completed.
func (srv *MockTransServer) Close() {
	_ = srv.listener.Close() // nolint: gosec
}

// Serve starts listening for requests
func (srv *MockTransServer) Serve(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go srv.handleRequest(conn)
	}
}

// handleRequest handles the request made to the server using the defined
// handler function
func (srv *MockTransServer) handleRequest(conn io.ReadWriteCloser) {
	defer conn.Close() // nolint: errcheck
	srv.mtx.RLock()
	busy := srv.IsBusy
	srv.mtx.RUnlock()

	if busy {
		_, err := conn.Write([]byte(BusyMessage))
		if err != nil {
			panic(err)
		}
		return
	}
	_, err := conn.Write([]byte(WelcomeMessage))
	if err != nil {
		panic(err)
	}

	br := bufio.NewReader(conn)
	var args []byte

	for {
		buf, err := br.ReadBytes('\n') // nolint: vetshadow
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}

		args = append(args, buf...)

		if bytes.Equal(buf, []byte(EndMessage)) {
			break
		}
	}

	// get the handler and pass the args to ger a response
	srv.mtx.RLock()
	h := srv.handler
	srv.mtx.RUnlock()

	if h != nil {
		res := h(args)
		_, err = conn.Write(res)
		if err != nil {
			panic(err)
		}
	}
	// add the end of the message
	_, err = conn.Write([]byte(EndMessage))
	if err != nil {
		panic(err)
	}
}

// SetHandler sets handler function.
func (srv *MockTransServer) SetHandler(h Handler) {
	srv.mtx.Lock()
	srv.handler = h
	srv.mtx.Unlock()
}

// SetBusy sets if the server is busy or available to response requests.
func (srv *MockTransServer) SetBusy(busy bool) {
	srv.mtx.Lock()
	srv.IsBusy = busy
	srv.mtx.Unlock()
}

// NewMockTransServer starts and returns a new Server.
// The caller should call Close when finished, to shut it down.
func NewMockTransServer() *MockTransServer {
	s := &MockTransServer{
		listener: newLocalListener(),
	}
	s.Start()
	return s
}

// newLocalListener starts a new TCP listener on the next available port
func newLocalListener() net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if l, err = net.Listen("tcp6", "[::1]:0"); err != nil {
			panic(fmt.Sprintf("transtest: failed to listen on a port: %v", err))
		}
	}
	return l
}
