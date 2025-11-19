package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/ABfry/album-battler/backend/internal/usecase/user"
)

type UserHandler struct {
	createUserUC *user.CreateUserUseCase
}

func NewUserHandler(createUserUC *user.CreateUserUseCase) *UserHandler {
	return &UserHandler{createUserUC: createUserUC}
}

// POST /user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserName string `json:"user_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	output, err := h.createUserUC.Execute(r.Context(), user.CreateUserInput{
		UserName: req.UserName,
	})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": output.UserID.String(),
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
