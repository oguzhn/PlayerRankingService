package main

import (
	"sync"

	"github.com/gorilla/handlers"
	"github.com/kelseyhightower/envconfig"

	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	const envPrefix = "RANKING"

	var conf struct {
		Mongo struct {
			Addr string `required:"true"`
			DB   string `required:"true"`
			Col  string `required:"true"`
		}
		Addr string `required:"true"`
	}

	err := envconfig.Process(envPrefix, &conf)
	if err != nil {
		envconfig.Usage(envPrefix, &conf)
		log.Fatalln("cannot parse env vars:", err)
	}
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	restController := controller.NewController()
	r := http.NewServeMux()
	r.Handle("/", restController.Handlers())

	server := &http.Server{
		Addr: conf.Addr,
		Handler: handlers.RecoveryHandler()(
			handlers.ProxyHeaders(
				handlers.LoggingHandler(os.Stdout,
					r,
				),
			),
		),
	}

	wg.Add(1)
	go func() {
		defer cancel()
		defer wg.Done()
		log.Println("Server started to listen")
		err := server.ListenAndServe()
		if err != nil {
			log.Println("server erred: ", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	select {
	case <-sigCh:
		cancel()
		log.Println("os signal caught ending")
	case <-ctx.Done():
		log.Println("context done")
	}

	sctx, sctxcancel := context.WithTimeout(context.Background(), time.Second)
	defer sctxcancel()
	err = server.Shutdown(sctx)
	if err != nil {
		log.Println("error closing server:", err)
	}
	wg.Wait()
}
