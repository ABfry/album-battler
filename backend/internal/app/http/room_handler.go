package httpapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/ABfry/album-battler/backend/internal/usecase/room"
	"github.com/google/uuid"
)

type RoomHandler struct {
	createRoomUC   *room.CreateRoomUseCase
	joinRoomUC     *room.JoinRoomUseCase
	startGameUC    *room.StartGameUseCase
	getRoomUC      *room.GetRoomUseCase
	roomManager    service.RoomManager
	eventPublisher service.EventPublisher
}

func NewRoomHandler(
	createRoomUC *room.CreateRoomUseCase,
	joinRoomUC *room.JoinRoomUseCase,
	startGameUC *room.StartGameUseCase,
	getRoomUC *room.GetRoomUseCase,
	roomManager service.RoomManager,
	eventPublisher service.EventPublisher,
) *RoomHandler {
	return &RoomHandler{
		createRoomUC:   createRoomUC,
		joinRoomUC:     joinRoomUC,
		startGameUC:    startGameUC,
		getRoomUC:      getRoomUC,
		roomManager:    roomManager,
		eventPublisher: eventPublisher,
	}
}

// POST /room
func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "Invalid user_id format", http.StatusBadRequest)
		return
	}

	output, err := h.createRoomUC.Execute(r.Context(), room.CreateRoomInput{
		UserID: userID,
	})
	if err != nil {
		log.Printf("CreateRoom error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"room_id":     output.RoomID.String(),
		"room_number": output.RoomNumber,
	})
}

// POST /room/join
func (h *RoomHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID     string `json:"user_id"`
		RoomNumber int    `json:"room_number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "Invalid user_id format", http.StatusBadRequest)
		return
	}

	err = h.joinRoomUC.Execute(r.Context(), room.JoinRoomInput{
		UserID:     userID,
		RoomNumber: req.RoomNumber,
	})
	if err != nil {
		log.Printf("JoinRoom error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: room_idを取得する(usecaseから返すようにする)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "joined successfully",
	})
}

// POST /room/{room_id}/start
func (h *RoomHandler) StartGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLからroom_idを取得 (簡易的にQuery Paramで受け取る)
	roomIDStr := r.URL.Query().Get("room_id")
	if roomIDStr == "" {
		http.Error(w, "room_id parameter is required", http.StatusBadRequest)
		return
	}

	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		http.Error(w, "Invalid room_id format", http.StatusBadRequest)
		return
	}

	var req struct {
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "Invalid user_id format", http.StatusBadRequest)
		return
	}

	err = h.startGameUC.Execute(r.Context(), room.StartGameInput{
		RoomID: roomID,
		UserID: userID,
	})
	if err != nil {
		log.Printf("StartGame error: %v", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"result": "ok",
	})
}

// GET /room/{room_id}
func (h *RoomHandler) GetRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLからroom_idを取得 (簡易的にQuery Paramで受け取る)
	roomIDStr := r.URL.Query().Get("room_id")
	if roomIDStr == "" {
		http.Error(w, "room_id parameter is required", http.StatusBadRequest)
		return
	}

	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		http.Error(w, "Invalid room_id format", http.StatusBadRequest)
		return
	}

	output, err := h.getRoomUC.Execute(r.Context(), room.GetRoomInput{
		RoomID: roomID,
	})
	if err != nil {
		log.Printf("GetRoom error: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
