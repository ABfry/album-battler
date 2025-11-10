package websocket

import "time"

// WebSocketで送受信するメッセージの型
type Message struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// ユーザーメッセージのペイロード
type UserMessagePayload struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}
