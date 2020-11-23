package main

import (
	"flag"
	"fmt"
	"github.com/francisco-serrano/sample-gokit/messages"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
	)

	flag.Parse()

	store := messages.NewInMemoryStore()
	svc := messages.NewMessageService(store)

	insertHandler := httptransport.NewServer(
		messages.MakeInsertEndpoint(svc),
		messages.DecodeInsertRequest,
		messages.EncodeResponse,
	)

	getAllHandler := httptransport.NewServer(
		messages.MakeGetAllEndpoint(svc),
		messages.DecodeGetAllRequest,
		messages.EncodeResponse,
	)

	r := mux.NewRouter()

	r.Methods(http.MethodPost).Path("/car").Handler(insertHandler)
	r.Methods(http.MethodGet).Path("/messages").Handler(getAllHandler)

	server := http.Server{
		Addr:    *listen,
		Handler: r,
	}

	errs := make(chan error, 2)
	go func() {
		errs <- server.ListenAndServe()
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	<-errs
}
