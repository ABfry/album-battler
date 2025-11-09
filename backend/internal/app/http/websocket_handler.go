package httpapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ABfry/album-battler/backend/internal/infra/websocket"
	"github.com/ABfry/album-battler/backend/internal/usecase"
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
	hub                     *websocket.Hub
	broadcastMessageUseCase *usecase.BroadcastMessageUseCase
}

func NewWebSocketHandler(
	hub *websocket.Hub,
	broadcastMessageUseCase *usecase.BroadcastMessageUseCase,
) *WebSocketHandler {
	return &WebSocketHandler{
		hub:                     hub,
		broadcastMessageUseCase: broadcastMessageUseCase,
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
	// WebSocket接続にアップグレード
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// 開発用: ランダムなユーザーIDを割り当て
	// 本番環境では認証トークンからユーザーIDを取得すべき
	userID := uuid.New()

	// クライアントを作成してHubに登録
	client := h.hub.NewClient(userID, conn)

	// goroutineでRead/Writeポンプを起動
	go client.WritePump()
	go client.ReadPump()

	log.Printf("New WebSocket connection: userID=%s", userID)
}
