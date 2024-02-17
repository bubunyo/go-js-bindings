package server

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	gojsbindings "github.com/bubunyo/go-js-bindings"
)

type Config struct {
	Dir string
}

type Server struct {
	fs fs.FS
}

func NewServer(config Config) *Server {
	s := &Server{
		// dir: config.Dir,
	}
	var err error
	s.fs, err = fs.Sub(gojsbindings.UI, "ui/build")
	if err != nil {
		log.Fatal("failed to get ui fs", err)
	}
	return s
}
func (s *Server) Stop() {
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	path := cleanPath(r.URL.Path)

	fmt.Println(">>>>>>>>>>>>", path)

	file, err := s.fs.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("file", path, "not found:", err)
			http.NotFound(w, r)
			return
		}
		log.Println("file", path, "cannot be read:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	contentType := mime.TypeByExtension(filepath.Ext(path))
	w.Header().Set("Content-Type", contentType)
	if strings.HasPrefix(path, "static/") {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
	}
	stat, err := file.Stat()
	if err == nil && stat.Size() > 0 {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	}

	n, _ := io.Copy(w, file)
	log.Println("file", path, "copied", n, "bytes")
}

func cleanPath(path string) string {
	path = filepath.Clean(path)
	path = strings.TrimPrefix(path, "/ui")
	if path == "/" || path == "" { // Add other paths that you route on the UI side here
		path = "index.html"
	}
	path = strings.TrimPrefix(path, "/")
	return path
}
