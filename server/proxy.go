package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

// Proxy - Manages a Proxy connection, piping data between local and remote.
type Proxy struct {
	conn     *net.TCPConn
	stopChan chan struct{}
}

// New - Create a new Proxy instance. Takes a remote address and creates
// a unix file socket. This allows you to manipulate the connection forwarding
// however you please
// and closes it when finished.
func NewProxy(conn *net.TCPConn) *Proxy {
	p := &Proxy{
		conn:     conn,
		stopChan: make(chan struct{}),
	}
	return p
}

func (p *Proxy) logError(msg string, err error) {
	log.Print("[Connection Proxy][Error]", msg, err.Error())
}

func (p *Proxy) logMsg(msg ...interface{}) {
	log.Println("[Connection Proxy]", msg)
}

func (p *Proxy) Stop() {
	p.logMsg("Intercepting and stopping connection")
	close(p.stopChan)
}

func (p *Proxy) start(w http.ResponseWriter, r *http.Request) {
	p.logMsg("new connection go go go")

	// p.stopChan = make(chan struct{})

	// select {
	// case <-p.stopChan:
	// 	goto exit
	// default:
	// }
	// lconn, err := listener.AcceptUnix()
	// if err != nil {
	// 	p.logError("Accept unix connection error", err)
	// 	return
	// }
	// p.logMsg("New connection:", lconn.RemoteAddr())

	// go func() {
	// 	defer func() {
	// 		// _ = lconn.Close()
	// 		// _ = listener.Close()
	// 	}()

	// connect to remote
	// rconn, err := net.DialTCP("tcp", nil, p.addr)
	// if err != nil {
	// 	p.logError("Remote connection failed:", err)
	// 	return
	// }
	// defer rconn.Close() //nolint:errcheck

	// display both ends
	// p.logMsg("Opened:: Remote:", p.addr.String())

	// closer := make(chan struct{})

	// reader, writer := io.Pipe()

	// go r.WriteProxy(writer)

	fmt.Fprintln(w, "hello world papi")

	return

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>. 4")

	// b := bytes.Buffer{}
	b := strings.Builder{}
	err := r.WriteProxy(&b)
	if err != nil {
		fmt.Println(">>>>>>>>>>>>>>>>> err", err)
	}
	fmt.Println(">>>>>>>>>>>>>>>>", b.String())

	// go func() {
	// 	if b, err := io.ReadAll(reader); err == nil {
	// 		fmt.Println(">>>>>>>>>>>>>>>>>>>> incoming", string(b))
	// 	}
	// }()
	// return
	// p.conn.Write()

	// p.conn.ReadFrom(reader)
	// p.conn.ReadFrom()
	// p.conn.Write()

	// bidirectional copy
	// go p.pipe(closer, reader, p.conn)
	// go p.pipe(closer, p.conn, w)

	// select {
	// case <-closer:
	// 	p.logMsg("Closing connection")
	// case <-p.stopChan:
	// 	p.logMsg("Stopping connection")
	// }

	p.logMsg("Connection complete")
	// }()
	// exit:
	// 	p.logMsg("Proxy shutting down")
}

func (p *Proxy) pipe(closer chan struct{}, src io.Reader, dst io.Writer) {
	_, _ = io.Copy(dst, src)
	closer <- struct{}{}
}
