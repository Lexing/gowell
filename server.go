package gowell

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"git.apache.org/thrift.git/lib/go/thrift"

	"github.com/gorilla/mux"
)

func defaultHealthzHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

type HttpServer struct {
	router         *mux.Router
	addr           string
	healthzHandler func(http.ResponseWriter, *http.Request)
}

func NewHttpServer(addr string) *HttpServer {
	s := &HttpServer{
		addr:           addr,
		router:         mux.NewRouter(),
		healthzHandler: defaultHealthzHandler,
	}

	return s
}

func (s *HttpServer) SetHealthzHandler(h func(http.ResponseWriter, *http.Request)) {
	s.healthzHandler = h
}

func (s *HttpServer) SetRouter(r *mux.Router) {
	s.router = r
}

func (s *HttpServer) SetAddr(addr string) {
	s.addr = addr
}

func (s *HttpServer) Start() {
	flag.Parse()
	s.router.HandleFunc("/healthz", s.healthzHandler)
	log.Printf("Http server ready to serve on %v", s.addr)
	err := http.ListenAndServe(s.addr, s.router)
	if err != nil {
		log.Panic(err)
	}
}

// ThriftServer wraps with HTTP server for basic monitoring.
type ThriftServer struct {
	// HTTP server for basic utils query
	http *HttpServer

	addr string
}

// NewThriftServer creates new ThriftServer listening on addr.
func NewThriftServer(addr string) *ThriftServer {
	s := &ThriftServer{
		addr: addr,
		http: NewHttpServer(":8080"),
	}

	return s
}

// Start starts Thrift server with given thrift processor
func (s *ThriftServer) Start(processor thrift.TProcessor) {
	flag.Parse()
	go s.http.Start()

	s.startThriftServer(processor)
}

func (s *ThriftServer) startThriftServer(processor thrift.TProcessor) {
	transport_factory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocol_factory := thrift.NewTCompactProtocolFactory()

	server_transport, err := thrift.NewTServerSocket(s.addr)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Thrift server listen on:", s.addr)
	server := thrift.NewTSimpleServer4(processor, server_transport, transport_factory, protocol_factory)
	err = server.Serve()
	if err != nil {
		log.Panic(err)
	}
}
