package httpapi

// -- ルーティングの定義 --

func (s *APIServer) registerRoutes() {
	// WebSocketハンドラ
	wsHandler := NewWebSocketHandler(
		s.deps.WebSocketHub)
	s.mux.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// ヘルスチェック
	// s.mux.HandleFunc("/healthz", s.handleHealthz)
}
