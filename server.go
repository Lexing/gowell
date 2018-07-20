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

var (
	port = flag.String("gowell_port", "8080", "gowell port")
)

var healthyCh = make(chan bool, 1)
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
	err := http.ListenAndServe(s.addr, s.router)
	if err != nil {
		log.Panic(err)
	}

	go func() {
		<-healthyCh
		hLock.Lock()
		healthy = true
		hLock.Unlock()
		log.Printf("Server is now healthy on %v.", s.addr)
	}()
}

// InitializeHTTPService starts a HTTP server and add basic http services, e.g. monitoring
func InitializeHTTPService() {
	s := NewHttpServer(fmt.Sprintf(":%v", *port))
	s.Start()
}

// NoteHealthy marks this server as healthy, reports 'ok' in /healthz
func NoteHealthy() {
	var once sync.Once
	once.Do(func() {
		healthyCh <- true
	})
}
