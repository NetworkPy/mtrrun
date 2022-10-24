package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/mtrrun/internal/handler"
	"github.com/mtrrun/internal/repository"
	"github.com/mtrrun/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// TODO: will be remove how project starts  to use config
var (
	defaultAddr = "127.0.0.1:8080"
)

func main() {
	r := mux.NewRouter()

	// Create repository layer
	metCache := repository.NewMetricMemCache()

	// Create service layer
	metSrv := service.NewMetricService(&service.MetricServiceConfig{
		MetRepo: metCache,
	})

	// Register all endpoints
	handler.New(&handler.Config{
		Router: r,
		MetSrv: metSrv,
	})

	srv := &http.Server{
		Addr:    defaultAddr,
		Handler: r,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Printf("server started on 8080 port")

	<-done
	log.Print("server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed:%+v", err)
	}
	log.Print("server exited properly")
}
