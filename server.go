package gowell

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO: add port flag
// var port = flag.String("util_port", "8080", "listening port for ")

var healthy bool
var hLock sync.RWMutex

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	hLock.RLock()
	defer hLock.RUnlock()
	if healthy {
		fmt.Fprint(w, "ok")
	}
}

func flagzHandler(w http.ResponseWriter, r *http.Request) {
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(w, "%v: %v\n", f.Name, f.Value)
	})
}

type HttpServer struct {
	router *mux.Router
	addr   string
}

func NewHttpServer(addr string) *HttpServer {
	s := &HttpServer{
		addr:   addr,
		router: mux.NewRouter(),
	}

	return s
}

func (s *HttpServer) SetRouter(r *mux.Router) {
	s.router = r
}

func (s *HttpServer) SetAddr(addr string) {
	s.addr = addr
}

func (s *HttpServer) Start() {
	flag.Parse()
	s.router.HandleFunc("/healthz", healthzHandler)
	s.router.HandleFunc("/flagz", flagzHandler)
	s.router.Handle("/metrics", promhttp.Handler())
	log.Printf("Http server ready to serve on %v", s.addr)
	err := http.ListenAndServe(s.addr, s.router)
	if err != nil {
		log.Panic(err)
	}
}

// InitializeHTTPService starts a HTTP server and add basic http services, e.g. monitoring
func InitializeHTTPService(addr string) {
	s := NewHttpServer(":8080")
	s.Start()
}

// NoteHealthy marks this server as healthy, reports 'ok' in /healthz
func NoteHealthy() {
	hLock.Lock()
	healthy = true
	hLock.Unlock()
}
