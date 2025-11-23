package httpapi

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ABfry/album-battler/backend/internal/domain/event"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
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

	roomRepo        repository.RoomRepository
	clapSendUseCase *clap.ClapSendUseCase
}

func NewWebSocketHandler(
	hub *websocket.Hub,
	roomRepo repository.RoomRepository,
	clapSendUseCase *clap.ClapSendUseCase,
) *WebSocketHandler {
	handler := &WebSocketHandler{
		hub:             hub,
		incoming:        make(chan *websocket.ClientInboundMessage, 256),
		roomRepo:        roomRepo,
		clapSendUseCase: clapSendUseCase,
	}
	go handler.consumeIncoming(context.Background())
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
	UserID       string `json:"user_id"`
	TargetUserID string `json:"target_user_id"`
	BattleID     string `json:"battle_id"`
	Count        int    `json:"count"`
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

	// ユーザーが参加中の部屋をDBから取得し、WebSocket Hub上のルームに再登録
	if h.roomRepo != nil {
		rooms, err := h.roomRepo.FindByUserID(r.Context(), userID)
		if err != nil {
			log.Printf("Failed to find rooms for user=%s: %v", userID, err)
		} else {
			for _, room := range rooms {
				h.hub.JoinRoom(userID, room.ID)
				log.Printf("Restored room membership: userID=%s, roomID=%s", userID, room.ID)
			}
		}
	}
}

func (h *WebSocketHandler) consumeIncoming(ctx context.Context) {
	for msg := range h.incoming {
		var incoming IncomingMessage
		if err := json.Unmarshal(msg.Payload, &incoming); err != nil {
			log.Printf("invalid incoming message from %s: %v", msg.UserID, err)
			continue
		}

		switch incoming.Type {
		case event.ClapSendEvent{}.EventType():
			h.handleClapSend(ctx, msg, incoming)
		default:
			log.Printf("unhandled ws message type=%s from user=%s", incoming.Type, msg.UserID)
		}
	}
}

func (h *WebSocketHandler) handleClapSend(ctx context.Context, msg *websocket.ClientInboundMessage, incoming IncomingMessage) {
	var payload ClapSendPayload
	if err := json.Unmarshal(incoming.Payload, &payload); err != nil {
		log.Printf("invalid clap_send payload from user=%s: %v", msg.UserID, err)
		return
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		log.Printf("invalid user_id in clap_send from user=%s: %v", msg.UserID, err)
		return
	}

	targetUserID, err := uuid.Parse(payload.TargetUserID)
	if err != nil {
		log.Printf("invalid target_user_id in clap_send from user=%s: %v", msg.UserID, err)
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

	if err := h.clapSendUseCase.Execute(ctx, clap.ClapSendInput{
		BattleID:     battleID,
		UserID:       userID,
		TargetUserID: targetUserID,
		Count:        payload.Count,
	}); err != nil {
		log.Printf("failed to execute ClapSendUseCase for user=%s: %v", msg.UserID, err)
	}
}
