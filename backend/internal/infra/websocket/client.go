package websocket

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// クライアントへの書き込み待機時間
	writeWait = 10 * time.Second

	// クライアントからのpong待機時間
	pongWait = 60 * time.Second

	// ping送信間隔（pongWaitより短くする必要がある）
	pingPeriod = (pongWait * 9) / 10

	// 最大メッセージサイズ
	maxMessageSize = 512
)

// Client は単一のWebSocket接続を表す
type Client struct {
	// クライアントのユーザーID
	UserID uuid.UUID

	// Hubへの参照
	Hub *Hub

	// WebSocket接続
	Conn *websocket.Conn

	// クライアントへの送信メッセージチャネル
	Send chan []byte

	// クライアントのルームID
	RoomID uuid.UUID
}

// 受信メッセージを読み取るループ
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Hubにメッセージを送信（ブロードキャスト用）
		c.Hub.broadcast <- &BroadcastMessage{
			SenderID: c.UserID,
			Message:  message,
		}
	}
}

// 書き込みループ, Sendから取り出す
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hubがチャネルを閉じた
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// キューにある追加メッセージも一緒に送信
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
