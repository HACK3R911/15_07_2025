package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"15_07_2025/config"
	"15_07_2025/pkg/handler"
	"15_07_2025/pkg/repositories"
	"15_07_2025/pkg/utils"
	"15_07_2025/pkg/worker"
)

func main() {
	port := flag.Int("port", 8080, "server port")
	extensions := flag.String("ext", "pdf,jpg,jpeg,png", "allowed extensions")
	flag.Parse()

	cfg := config.NewConfig(*port, utils.ParseExtensions(*extensions))
	repository := repositories.NewRepository()
	wp := worker.NewWorkerPool(3, repository)
	handler := handler.NewArchiveHandler(repository, wp, cfg)

	wp.Start()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: handler,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server started on port %d", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-done
	log.Println("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)

	wp.Stop()
	log.Println("Server exited.")
}
