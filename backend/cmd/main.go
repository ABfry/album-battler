package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ABfry/album-battler/backend/internal/app"
	httpapi "github.com/ABfry/album-battler/backend/internal/app/http"
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
	go deps.WebSocketHub.Run(ctx) // Hubをgoroutineで実行

	// テスト用定期ブロードキャストを開始（10秒ごとに"Hello"を送信）
	// go deps.PeriodicBroadcastUseCase.Start(ctx)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := api.ListenAndServe(ctx, ":"+port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
