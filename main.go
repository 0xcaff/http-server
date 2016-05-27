package main

import (
  "net/http"
  "github.com/koding/logging"
  "flag"
  "fmt"
)

var (
  port = flag.String("port", "8080", "Port to use")
  path = flag.String("path", "./", "Path the server is started from")
)

func main() {
  flag.Parse()

  server := HTTPServer{
    http.Server{
      Addr: ":" + *port,
      Handler: HTTPHandler{
        realHandler: http.FileServer(http.Dir(*path)),
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
}

func (h HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  logging.Info("%s %s", r.Method, r.URL)

  h.realHandler.ServeHTTP(w, r)
}

type HTTPServer struct {
  http.Server
}

func (h *HTTPServer) ListenAndServe() error {
  logging.Info(
    "Starting server on \033[%dm%s\033[0m at \033[%dm%s\033[0m",
    logging.CYAN,
    h.Server.Addr,
    logging.CYAN,
    *path,
  )

  return h.Server.ListenAndServe()
}

