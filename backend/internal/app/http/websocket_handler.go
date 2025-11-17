package httpapi

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ABfry/album-battler/backend/internal/infra/websocket"
	"github.com/ABfry/album-battler/backend/internal/usecase/clap"
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

	clapSendUseCase *clap.ClapSendUseCase
}

func NewWebSocketHandler(
	hub *websocket.Hub,
	clapSendUseCase *clap.ClapSendUseCase,
) *WebSocketHandler {
	handler := &WebSocketHandler{
		hub:             hub,
		incoming:        make(chan *websocket.ClientInboundMessage, 256),
		clapSendUseCase: clapSendUseCase,
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

type ClapSendPayload struct {
	BattleID string `json:"battle_id"`
	Count    int    `json:"count"`
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

		switch incoming.Type {
		case "clap_send":
			h.handleClapSend(msg, incoming)
		default:
			log.Printf("unhandled ws message type=%s from user=%s", incoming.Type, msg.UserID)
		}
	}
}

func (h *WebSocketHandler) handleClapSend(msg *websocket.ClientInboundMessage, incoming IncomingMessage) {
	var payload ClapSendPayload
	if err := json.Unmarshal(incoming.Payload, &payload); err != nil {
		log.Printf("invalid clap_send payload from user=%s: %v", msg.UserID, err)
		return
	}

	battleID, err := uuid.Parse(payload.BattleID)
	if err != nil {
		log.Printf("invalid battle_id in clap_send from user=%s: %v", msg.UserID, err)
		return
	}

	if h.clapSendUseCase == nil {
		log.Printf("ClapSendUseCase is not initialized")
		return
	}

	if err := h.clapSendUseCase.Execute(context.Background(), clap.ClapSendInput{
		BattleID: battleID,
		UserID:   msg.UserID,
		Count:    payload.Count,
	}); err != nil {
		log.Printf("failed to execute ClapSendUseCase for user=%s: %v", msg.UserID, err)
	}
}
