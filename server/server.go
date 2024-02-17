package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"golang.org/x/exp/slog"
)

type Config struct {
	Dir string
}

type Server struct {
	dir     string
	process *os.Process
	port    int
	proxy   *Proxy
}

func NewServer(config Config) *Server {
	s := &Server{
		dir: config.Dir,
	}
	slog.Info("running command")
	return s
}

func (s *Server) Start() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// p, err := GetFreePort()
		// if err != nil {
		// 	fmt.Println(">>>>>>>>>>>>>>>>> port", err)
		// 	return
		// }
		p := 3000
		s.port = p
		log.Println("got free port", s.port)
		cmd := exec.Command("npm", "start")
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, fmt.Sprintf("PORT=%d BROWSER=none PUBLIC_URL='http:localhost:3000/ui/'", p))
		cmd.Dir = s.dir
		pipe, _ := cmd.StdoutPipe()
		errPipe, _ := cmd.StderrPipe()
		if err := cmd.Start(); err != nil {
			fmt.Println(">>>>>>>>>>>>>>>>> start", err)
			return
		}

		s.process = cmd.Process

		go logOut(pipe)
		go logOut(errPipe)

		wg.Done()

		if err := cmd.Wait(); err != nil {
			// handle error
			fmt.Println(">>>>>>>>>>>>>>>>>wait ", err)
			return
		}
	}()
	wg.Wait()
	slog.Info("front end started successfully", "port", s.port, "pid", s.process.Pid)
	time.Sleep(3 * time.Second)
	slog.Info("Starting Proxy")
	s.startProxy()
}

func (s *Server) Stop() {
	// if err := s.proxy.conn.Close(); err != nil {
	// 	slog.Info(">>> failed to close connection", "folder", s.dir, "pid", s.process.Pid, "port", s.port)
	// }

	if err := s.process.Kill(); err != nil {
		slog.Info(">>> failed to kill server running", "folder", s.dir, "pid", s.process.Pid, "port", s.port)
	}
	if err := s.process.Release(); err != nil {
		slog.Info(">>> failed to release process", "folder", s.dir, "pid", s.process.Pid, "port", s.port)
	}
}

func (s *Server) startProxy() {
	slog.Info(">>>>>>>> starting proxy")
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", s.port))
	if err != nil {
		slog.Error("Resolve remote address error", "error", err)
		return
	}

	// ln, err := net.ListenTCP("tcp", addr)
	// if err != nil {
	// 	slog.Error("listen failed", "error", err)
	// 	return
	// }
	// s.listener = ln

	rconn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		slog.Error("front end conn failed", "error", err)
		return
	}
	s.proxy = NewProxy(rconn)
	slog.Error(">>>>>>>>>>>>>>>>>>>>Proxy Started")
}

func logOut(pipe io.ReadCloser) {
	reader := bufio.NewReader(pipe)
	line, err := reader.ReadString('\n')
	for err == nil {
		fmt.Print(line)
		line, err = reader.ReadString('\n')
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>>>>>>>>>>>>>>.. print url", r.URL.String())
	s.proxy.start(w, r)
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
