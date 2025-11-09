package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ABfry/album-battler/backend/internal/app"
	httpapi "github.com/ABfry/album-battler/backend/internal/app/http"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	deps, err := app.NewDependencies()
	if err != nil {
		log.Fatalf("init dependencies: %v", err)
	}
	defer func() {
		if cerr := deps.Close(); cerr != nil {
			log.Printf("close dependencies: %v", cerr)
		}
	}()

	api := httpapi.NewAPIServer(deps)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := api.ListenAndServe(ctx, ":"+port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
