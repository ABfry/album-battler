package httpapi

// -- ルーティングの定義 --

func (s *APIServer) registerRoutes() {
	// WebSocketハンドラ
	wsHandler := NewWebSocketHandler(
		s.deps.WebSocketHub,
		s.deps.ClapSendUseCase,
	)
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

	// Battleハンドラ
	battleHandler := NewBattleHandler(
		s.deps.CreateBattleUseCase,
		s.deps.GetBattleUseCase,
		s.deps.GetBattleIDUseCase,
		s.deps.GetActiveBattleUseCase,
		s.deps.GetImageUseCase,
		s.deps.ImageSendUseCase,
		s.deps.GetResultUseCase,
	)

	// Userハンドラ
	userHandler := NewUserHandler(s.deps.CreateUserUseCase)
	s.mux.HandleFunc("POST /user", userHandler.CreateUser)

	s.mux.HandleFunc("POST /battle", battleHandler.CreateBattle) // デバッグ用
	s.mux.HandleFunc("GET /battle/{id}", battleHandler.GetBattle)
	s.mux.HandleFunc("GET /room/{id}/battle-id", battleHandler.GetBattleIDByRoom)
	s.mux.HandleFunc("GET /battle/{id}/image", battleHandler.GetImage)
	s.mux.HandleFunc("POST /battle/{id}/send-image", battleHandler.SendImage)
	s.mux.HandleFunc("GET /battle/{id}/result", battleHandler.GetResult)

	// ヘルスチェック
	s.mux.HandleFunc("GET /health", s.handleHealth)
}
