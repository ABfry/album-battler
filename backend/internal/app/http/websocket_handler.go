package httpapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ABfry/album-battler/backend/internal/infra/websocket"
	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CORS設定
	CheckOrigin: func(r *http.Request) bool {
		return true // 開発用: すべてのオリジンを許可
	},
}

type WebSocketHandler struct {
	hub *websocket.Hub
}

func NewWebSocketHandler(
	hub *websocket.Hub,
) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

// クライアントから受信するメッセージ
type IncomingMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// メッセージペイロード
type MessagePayload struct {
	Message string `json:"message"`
}

// WebSocket接続を処理
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// URLパラメータからUserIDを取得
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id format", http.StatusBadRequest)
		return
	}

	// WebSocket接続にアップグレード
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// クライアントを作成してHubに登録
	client := h.hub.NewClient(userID, conn)

	// goroutineでRead/Writeポンプを起動
	go client.WritePump()
	go client.ReadPump()

	log.Printf("New WebSocket connection: userID=%s", userID)
}
