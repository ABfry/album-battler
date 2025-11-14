package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/ABfry/album-battler/backend/internal/usecase/battle"
	"github.com/google/uuid"
)

type BattleHandler struct {
	getBattleUC *battle.GetBattleUseCase
}

func NewBattleHandler(
	getBattleUC *battle.GetBattleUseCase,
) *BattleHandler {
	return &BattleHandler{
		getBattleUC: getBattleUC,
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
