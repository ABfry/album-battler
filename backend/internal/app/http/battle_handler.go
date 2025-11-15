package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ABfry/album-battler/backend/internal/usecase/battle"
	"github.com/google/uuid"
)

type BattleHandler struct {
	createBattleUC   *battle.CreateBattleUseCase
	getBattleUC      *battle.GetBattleUseCase
	getBattleIDUseUC *battle.GetBattleIDUseCase
	getImageSendUC   *battle.ImageSendUseCase
}

func NewBattleHandler(
	createBattleUC *battle.CreateBattleUseCase,
	getBattleUC *battle.GetBattleUseCase,
	getBattleIDUC *battle.GetBattleIDUseCase,
	imageSendUC *battle.ImageSendUseCase,
) *BattleHandler {
	return &BattleHandler{
		createBattleUC:   createBattleUC,
		getBattleUC:      getBattleUC,
		getBattleIDUseUC: getBattleIDUC,
		getImageSendUC:   imageSendUC,
	}
}

// POST /battle
// デバッグ用エンドポイント 実際はStartRoom時にbattleを作成したい
func (h *BattleHandler) CreateBattle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		RoomID string `json:"room_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	roomID, err := uuid.Parse(req.RoomID)
	if err != nil {
		http.Error(w, "Invalid room_id format", http.StatusBadRequest)
		return
	}

	battleID, err := h.createBattleUC.Execute(r.Context(), battle.CreateBattleInput{
		RoomID: roomID,
	})
	if err != nil {
		http.Error(w, "Failed to create battle", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(battleID); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GET /battle/{id}
func (h *BattleHandler) GetBattle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	battleIDStr := r.PathValue("id")
	battleID, err := uuid.Parse(battleIDStr)
	if err != nil {
		http.Error(w, "Invalid battle_id format", http.StatusBadRequest)
		return
	}

	battle, err := h.getBattleUC.Execute(r.Context(), battle.GetBattleInput{
		BattleID: battleID,
	})
	if err != nil {
		fmt.Println("failed to get battle", err)
		http.Error(w, "Failed to get battle", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(battle); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GET /room/{id}/battle-id
func (h *BattleHandler) GetBattleIDByRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	roomIDStr := r.PathValue("id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		http.Error(w, "Invalid room_id format", http.StatusBadRequest)
		return
	}

	result, err := h.getBattleIDUseUC.Execute(r.Context(), battle.GetBattleIDInput{
		RoomID: roomID,
	})
	if err != nil {
		http.Error(w, "Failed to get battle id", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// POST /battle/{id}/send-image
func (h *BattleHandler) SendImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	battleIDStr := r.PathValue("id")
	battleID, err := uuid.Parse(battleIDStr)
	if err != nil {
		http.Error(w, "Invalid battle_id format", http.StatusBadRequest)
		return
	}

	var req struct {
		UserID      string `json:"user_id"`
		ImageBase64 []byte `json:"image_base64"`
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

	// コンテンツタイプを検出
	contentType := http.DetectContentType(req.ImageBase64)

	err = h.getImageSendUC.Execute(r.Context(), battle.ImageSendInput{
		UserID:      userID,
		BattleID:    battleID,
		ImageBase64: req.ImageBase64,
		ContentType: contentType,
	})
	if err != nil {
		http.Error(w, "Failed to send image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"result": "ok",
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
