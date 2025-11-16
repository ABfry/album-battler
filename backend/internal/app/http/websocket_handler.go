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
	hub      *websocket.Hub
	incoming chan *websocket.ClientInboundMessage
}

func NewWebSocketHandler(
	hub *websocket.Hub,
) *WebSocketHandler {
	handler := &WebSocketHandler{
		hub:      hub,
		incoming: make(chan *websocket.ClientInboundMessage, 256),
	}
	go handler.consumeIncoming()
	return handler
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
	client := h.hub.NewClient(userID, conn, h.incoming)

	// goroutineでRead/Writeポンプを起動
	go client.WritePump()
	go client.ReadPump()

	log.Printf("New WebSocket connection: userID=%s", userID)
}

func (h *WebSocketHandler) consumeIncoming() {
	for msg := range h.incoming {
		var incoming IncomingMessage
		if err := json.Unmarshal(msg.Payload, &incoming); err != nil {
			log.Printf("invalid incoming message from %s: %v", msg.UserID, err)
			continue
		}

		// TODO: ここでClapなどのメッセージ種別に応じたUseCaseを呼び出す
		switch incoming.Type {
		default:
			log.Printf("unhandled ws message type=%s from user=%s", incoming.Type, msg.UserID)
		}
	}
}
