package httpapi

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ABfry/album-battler/backend/internal/app"
)

type APIServer struct {
	deps       *app.Dependencies
	mux        *http.ServeMux
	httpServer *http.Server
}

func NewAPIServer(deps *app.Dependencies) *APIServer {
	srv := &APIServer{
		deps: deps,
		mux:  http.NewServeMux(),
	}
	srv.registerRoutes()
	return srv
}

// サーバーの起動
func (s *APIServer) ListenAndServe(ctx context.Context, addr string) error {
	s.httpServer = &http.Server{
		Addr:              addr,
		Handler:           s.withCORS(s.mux),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			// todo : シャットダウン時エラー
			fmt.Println("シャットダウン時エラー", err)
		}
	}()

	err := s.httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// CORSの設定
func (s *APIServer) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ヘルスチェック
func (s *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
