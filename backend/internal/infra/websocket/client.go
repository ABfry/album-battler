package websocket

import (
	"log"
	"sync"
	"sync/atomic"
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

	// close管理
	closeDone sync.Once
	closed    atomic.Bool // クライアントが閉じられたかどうか
}

// 受信メッセージを読み取るループ
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		if err := c.Conn.Close(); err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}()

	if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("failed to set read deadline: %v", err)
		return
	}
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetPongHandler(func(string) error {
		if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			log.Printf("failed to set read deadline in pong handler: %v", err)
		}
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
		if err := c.Conn.Close(); err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("failed to set write deadline: %v", err)
				return
			}
			if !ok {
				// Hubがチャネルを閉じた
				if err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Printf("failed to write close message: %v", err)
				}
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err := w.Write(message); err != nil {
				log.Printf("failed to write message: %v", err)
				return
			}

			// キューにある追加メッセージも一緒に送信
			n := len(c.Send)
			for i := 0; i < n; i++ {
				if _, err := w.Write([]byte{'\n'}); err != nil {
					log.Printf("failed to write newline: %v", err)
					return
				}
				if _, err := w.Write(<-c.Send); err != nil {
					log.Printf("failed to write queued message: %v", err)
					return
				}
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("failed to set write deadline for ping: %v", err)
				return
			}
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) IsClosed() bool {
	return c.closed.Load()
}

func (c *Client) Close() {
	c.closeDone.Do(func() {
		c.closed.Store(true)
		close(c.Send)
		if c.Conn != nil {
			c.Conn.Close()
		}
	})
}
