package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/koding/logging"
)

var (
	listen    = flag.String("listen", ":8080", "Address to listen for requests on")
	servePath = flag.String("path", "./", "Path the server is started from")
	header    = flag.String("header", "", "Header sent with every response")
	redirect  = flag.String("redirect", "", "Respond to unknown requests with this file. Used to serve Single Page Applications.")

	proxyFrom = flag.String("proxy-from", "", "Proxy requests to this path to the the path specified by proxy-to.")
	proxyTo   = flag.String("proxy-to", "", "Proxy requests from proxy-from to this address")
)

func main() {
	flag.Parse()

	logging.Info("Starting server on %s at %s", emp(*listen), emp(*servePath))

	var headers map[string]string
	if *header != "" {
		s := strings.Split(*header, ":")
		headers = map[string]string{
			s[0]: s[1],
		}
	}

	for k, v := range headers {
		logging.Info("Adding Header: %s:%s", emp(k), emp(v))
	}

	mux := http.NewServeMux()

	if *redirect != "" {
		logging.Info("Redirecting 404s to: %s", emp(*redirect))
	}

	mux.Handle("/", HeaderHandler{
		Handler: http.FileServer(SinglePageFileSystem{
			backendSystem: http.Dir(*servePath),
			redirectTo:    *redirect,
		}),
		headers: headers,
	})

	if *proxyFrom != "" && *proxyTo != "" {
		// Handle Redirects
		logging.Info("Proxying: (%s) -> (%s)",
			emp(*proxyFrom),
			emp(*proxyTo),
		)

		proxyToUrl, err := url.Parse(*proxyTo)
		if err != nil {
			logging.Error("Bad -proxy-to URL:", err)
			return
		}

		basePath, err := url.Parse(*proxyFrom)
		if err != nil {
			logging.Error("Bad -proxy-from URL:", err)
			return
		}

		director := httputil.NewSingleHostReverseProxy(proxyToUrl).Director
		proxy := httputil.ReverseProxy{
			Director: func(r *http.Request) {
				// Remove Base Path
				oldPath := *r.URL
				r.URL.Path = strings.TrimPrefix(oldPath.Path, basePath.Path)

				logging.Info("[proxying] (%s) -> (%s)", emp(oldPath.String()), emp(r.URL.String()))
				director(r)
			},
			Transport: &ProxyRoundTripper{
				BasePath: basePath.Path,
			},
		}

		mux.Handle(*proxyFrom, &proxy)
	}

	server := http.Server{
		Addr:    *listen,
		Handler: LoggingHandler{mux},
	}

	err := server.ListenAndServe()
	if err != nil {
		logging.Error(fmt.Sprint(err))
	}
}

type HeaderHandler struct {
	headers map[string]string
	http.Handler
}

func (h HeaderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	responseHeaders := w.Header()
	for k, v := range h.headers {
		responseHeaders.Set(k, v)
	}

	h.Handler.ServeHTTP(w, r)
}

type LoggingHandler struct {
	http.Handler
}

func (h LoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logging.Info("%s %s", r.Method, r.URL)

	h.Handler.ServeHTTP(w, r)
}

func colorize(color logging.Color, s string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", color, s)
}

func emp(s string) string {
	return colorize(logging.CYAN, s)
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
		logging.Info("[redirecting] (%s) -> (%s)\n", emp(name), emp(localPath))
	} else {
		localPath = name
	}

	return spa.backendSystem.Open(localPath)
}

type ProxyRoundTripper struct {
	BasePath string
}

func (t *ProxyRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	// Round Trip
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	// Rewrite Cookie Paths
	cookies := resp.Cookies()
	for i := range cookies {
		cookie := cookies[i]
		cookie.Path = path.Join(t.BasePath, cookie.Path)
	}

	// Remove Cookies
	resp.Header.Del("Set-Cookie")

	// Set Cookies
	for _, cookie := range cookies {
		resp.Header.Add("Set-Cookie", cookie.String())
	}

	return resp, nil
}
