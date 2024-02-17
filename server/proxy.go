package server

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Proxy - Manages a Proxy connection, piping data between local and remote.
type Proxy struct {
	// conn     *net.TCPConn
	// stopChan chan struct{}
	proxy *httputil.ReverseProxy
}

// New - Create a new Proxy instance. Takes a remote address and creates
// a unix file socket. This allows you to manipulate the connection forwarding
// however you please
// and closes it when finished.
func NewProxy(conn *net.TCPConn) *Proxy {
	// httpTransport := httpclient.NewTransport(
	// 	// hard coded to simplify the setup and testing; we can move it to the config later
	// 	httpclient.WithMaxIdleConnsPerHost(100),
	// )

	// director := func(r *http.Request) {
	// 	// pass proxy's home region for debugging

	// 	// if host := r.Host; host != "" {
	// 	// 	// for client requests, Request.Host overrides the Host header to send;
	// 	// 	// we never want to proxy Host header to avoid confusion, when observe http.Request.Host inside DR instances;
	// 	// 	// i.e. adjust_server sets host to bagger's logEntry for debugging; w/o removing the header,
	// 	// 	// the value contains the original request's host (refer to https://pkg.go.dev/net/http#Request.Host)
	// 	// 	r.Host = ""
	// 	// 	r.Header.Set(HttpHeaderAdjustDrProxyHost, host)
	// 	// }

	// 	// if _, ok := r.Header["User-Agent"]; !ok {
	// 	// 	// explicitly disable User-Agent so it's not set to default value
	// 	// 	r.Header.Set("User-Agent", "")
	// 	// }
	// }
	// errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
	// 	// metricsBuffer.CountM("proxy.failed.redirect")
	// 	// log.Error(
	// 	// 	ctx, "[DR-Middleware]: failed to proxy request, redirecting",
	// 	// 	log.Field("url", r.URL.String()), log.WithError(err),
	// 	// )
	// 	// http.Redirect(w, r, r.URL.String(), http.StatusFound)
	// 	fmt.Fprintln(w, "error handle going")
	// }
	// proxy := &httputil.ReverseProxy{
	// 	Director:     director,
	// 	ErrorHandler: errorHandler,
	// 	// Transport:    transport,
	// }
	// return &ProxyHandler{
	// 	proxy:         proxy,
	// 	handler:       handler,
	// }
	remote, err := url.Parse("http://localhost:3000")
	if err != nil {
		panic(err)
	}

	p := &Proxy{
		// conn:     conn,
		// stopChan: make(chan struct{}),
		proxy: httputil.NewSingleHostReverseProxy(remote),
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
	// close(p.stopChan)
}

func (p *Proxy) start(w http.ResponseWriter, r *http.Request) {
	p.logMsg("new connection go go go")
	log.Println(r.URL)
	w.Header().Set("X-Ben", "Rad")
	p.proxy.ServeHTTP(w, r)

	// remote, err := url.Parse("http://google.com")
	// if err != nil {
	// 	panic(err)
	// }
	// proxy := httputil.NewSingleHostReverseProxy(remote)

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

	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>. 4")

	// b := bytes.Buffer{}
	// // b := strings.Builder{}
	// err := r.WriteProxy(&b)
	// if err != nil {
	// 	log.Println(">>>>>>>>>>>>>>>>> err", err)
	// }
	// log.Println(">>>>>>>>>>>>>>>>", b.String())

	// // go func() {
	// // 	if b, err := io.ReadAll(reader); err == nil {
	// // 		fmt.Println(">>>>>>>>>>>>>>>>>>>> incoming", string(b))
	// // 	}
	// // }()
	// // return
	// bb, _ := b.ReadBytes(0)
	// // p.conn.Write(bb)

	// // p.conn.ReadFrom(reader)
	// // p.conn.ReadFrom()
	// // p.conn.Write()

	// closer := make(chan struct{})

	// slog.Info(">>>>>>>>>>>>>>> piping")

	// // bidirectional copy
	// go p.pipe(closer, bytes.NewReader(bb), p.conn)
	// go p.pipe(closer, p.conn, w)

	// select {
	// case <-closer:
	// 	p.logMsg("Closing connection")
	// case <-p.stopChan:
	// 	p.logMsg("Stopping connection")
	// }

	// p.logMsg("Connection complete")
	// }()
	// exit:
	// 	p.logMsg("Proxy shutting down")

	// fmt.Fprintln(w, "hello world papi")
	// url, _ := url.Parse("http://localhost:3000")

	// create the reverse proxy
	// proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	// r.URL.Host = url.Host
	// r.URL.Scheme = url.Scheme
	// r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	// r.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	// proxy.ServeHTTP(w, r)
	// p.proxy.ServeHTTP(w, r)

	return
}

// func (p *Proxy) pipe(closer chan struct{}, src io.Reader, dst io.Writer) {
// 	_, _ = io.Copy(dst, src)
// 	// closer <- struct{}{}
// }
