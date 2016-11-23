package main

import (
	"flag"
	"fmt"
	"github.com/koding/logging"
	"net/http"
	"strings"
)

var (
	listen = flag.String("listen", ":8080", "Address to listen for requests on")
	path   = flag.String("path", "./", "Path the server is started from")
	header = flag.String("header", "", "Header sent with every response")
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
				realHandler: http.FileServer(http.Dir(*path)),
				headers:     headers,
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
