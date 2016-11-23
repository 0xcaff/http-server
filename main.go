package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/koding/logging"
)

var (
	listen   = flag.String("listen", ":8080", "Address to listen for requests on")
	path     = flag.String("path", "./", "Path the server is started from")
	header   = flag.String("header", "", "Header sent with every response")
	redirect = flag.String("redirect", "", "Respond to unknown requests with this file. Used to serve Single Page Applications.")
)

func main() {
	flag.Parse()

	var headers map[string]string
	if *header != "" {
		s := strings.Split(*header, ":")
		headers = map[string]string{
			s[0]: s[1],
		}
	}

	server := HTTPServer{
		http.Server{
			Addr: *listen,
			Handler: HTTPHandler{
				realHandler: http.FileServer(SinglePageFileSystem{
					backendSystem: http.Dir(*path),
					redirectTo:    *redirect,
				}),
				headers: headers,
			},
		},
	}

	err := server.ListenAndServe()
	if err != nil {
		logging.Error(fmt.Sprint(err))
	}
}

type HTTPHandler struct {
	realHandler http.Handler
	headers     map[string]string
}

func (h HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logging.Info("%s %s", r.Method, r.URL)

	responseHeaders := w.Header()
	for k, v := range h.headers {
		responseHeaders.Set(k, v)
	}

	h.realHandler.ServeHTTP(w, r)
}

type HTTPServer struct {
	http.Server
}

func (h *HTTPServer) ListenAndServe() error {
	logging.Info(
		"Starting server on %s at %s",
		colorize(logging.CYAN, h.Server.Addr),
		colorize(logging.CYAN, *path),
	)

	for k, v := range h.Handler.(HTTPHandler).headers {
		logging.Info("Adding Header: %s:%s",
			colorize(logging.CYAN, k),
			colorize(logging.CYAN, v),
		)
	}

	return h.Server.ListenAndServe()
}

func colorize(color logging.Color, s string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", color, s)
}

// A http.FileSystem for serving single page applications by redirecting all
// unknown paths to a given path.
type SinglePageFileSystem struct {
	backendSystem http.Dir
	redirectTo    string
}

func (spa SinglePageFileSystem) Open(name string) (http.File, error) {
	var localPath string
	basePath := string(spa.backendSystem)
	reqPath := filepath.Join(basePath, name)

	if _, err := os.Stat(reqPath); os.IsNotExist(err) && spa.redirectTo != "" {
		localPath = spa.redirectTo
		logging.Info("[redirecting] (%s) -> (%s)\n",
			colorize(logging.CYAN, name),
			colorize(logging.CYAN, localPath),
		)
	} else {
		localPath = name
	}

	return spa.backendSystem.Open(localPath)
}
