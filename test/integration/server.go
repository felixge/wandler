package integration

import (
	"net"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Get name/dir of this source file
var (
	_, __filename, _, _ = runtime.Caller(0)
	__dirname           = filepath.Dir(__filename)
)

func NewServer() (*Server, error) {
	p, err := newServerProcess()
	if err != nil {
		return nil, err
	}
	server := &Server{process: p}
	if err := server.waitUntilReady(); err != nil {
		server.Kill()
		return nil, err
	}
	return server, nil
}

func newServerProcess() (*os.Process, error) {
	bin := "wandler-server"
	path := filepath.Join(__dirname, "../../bin/", bin)
	return os.StartProcess(path, []string{bin}, &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	})
}

type Server struct {
	process *os.Process
}

func (s *Server) Kill() error {
	s.process.Signal(os.Interrupt)
	_, err := s.process.Wait()
	return err
}

// @TODO don't hardcode this, pass a port as a flag when creating the server
// and return this port here.
func (s *Server) HttpAddr() string {
	return "localhost:8080"
}

func (s *Server) waitUntilReady() error {
	for {
		conn, err := net.Dial("tcp", s.HttpAddr())
		if err == nil {
			return conn.Close()
		}
		time.Sleep(100 * time.Millisecond)
	}
}
