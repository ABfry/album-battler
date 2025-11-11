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
		s.deps.StartGameUseCase,
		s.deps.GetRoomUseCase,
		s.deps.RoomManager,
		s.deps.EventPublisher,
	)
	s.mux.HandleFunc("/room", roomHandler.CreateRoom)
	s.mux.HandleFunc("/room/join", roomHandler.JoinRoom)
	s.mux.HandleFunc("/room/start", roomHandler.StartGame)
	s.mux.HandleFunc("/room/info", roomHandler.GetRoom)

	// ヘルスチェック
	// s.mux.HandleFunc("/healthz", s.handleHealthz)
}
