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
	roomManager service.RoomManager,
	eventPublisher service.EventPublisher,
) *RoomHandler {
	return &RoomHandler{
		createRoomUC:   createRoomUC,
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
