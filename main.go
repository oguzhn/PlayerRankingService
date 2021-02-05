package main

import (
	"sync"

	"github.com/gorilla/handlers"
	"github.com/kelseyhightower/envconfig"
	"github.com/oguzhn/PlayerRankingService/business"
	"github.com/oguzhn/PlayerRankingService/controller"
	"github.com/oguzhn/PlayerRankingService/database"

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
			DB   string `default:"RankingDb"`
			Col  string `default:"Players"`
		}
		Addr string `default:":8080"`
	}

	err := envconfig.Process(envPrefix, &conf)
	if err != nil {
		envconfig.Usage(envPrefix, &conf)
		log.Fatalln("cannot parse env vars:", err)
	}
	mongoclient, err := database.NewDatastore(conf.Mongo.Addr, conf.Mongo.DB, conf.Mongo.Col)
	if err != nil {
		log.Fatalln("cannot reach mongodb", err)
	}
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	businessLayer := business.NewBusiness(mongoclient)

	restController := controller.NewController(businessLayer)
	r := http.NewServeMux()
	r.Handle("/", restController.RegisterHandlers())

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
