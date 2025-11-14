package httpapi

// -- ルーティングの定義 --

func (s *APIServer) registerRoutes() {
	// WebSocketハンドラ
	wsHandler := NewWebSocketHandler(
		s.deps.WebSocketHub)
	s.mux.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// Roomハンドラ
	roomHandler := NewRoomHandler(
		s.deps.CreateRoomUseCase,
		s.deps.JoinRoomUseCase,
		s.deps.LeaveRoomUseCase,
		s.deps.StartGameUseCase,
		s.deps.GetRoomUseCase,
		s.deps.RoomManager,
		s.deps.EventPublisher,
	)
	s.mux.HandleFunc("POST /room", roomHandler.CreateRoom)
	s.mux.HandleFunc("POST /room/join", roomHandler.JoinRoom)
	s.mux.HandleFunc("POST /room/{id}/leave", roomHandler.LeaveRoom)
	s.mux.HandleFunc("POST /room/{id}/start", roomHandler.StartGame)
	s.mux.HandleFunc("GET /room/{id}", roomHandler.GetRoom)

	// ヘルスチェック
	// s.mux.HandleFunc("/healthz", s.handleHealthz)
}
