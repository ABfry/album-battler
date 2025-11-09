package websocket

import (
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// ブロードキャストするメッセージ
type BroadcastMessage struct {
	SenderID uuid.UUID // 送信者のユーザーID
	Message  []byte
}

// 部屋単位のブロードキャストメッセージ
type RoomBroadcastMessage struct {
	RoomID  uuid.UUID
	Message interface{}
}

// 特定ユーザーへのメッセージ
type UserMessage struct {
	UserID  uuid.UUID
	Message interface{}
}

// アクティブなクライアントを管理
type Hub struct {
	// 登録されたクライアント (ユーザーID → Client)
	clients map[uuid.UUID]*Client

	// 部屋のメンバー管理 (部屋ID → ユーザーIDのセット)
	rooms map[uuid.UUID]map[uuid.UUID]bool

	// クライアントからのブロードキャストメッセージ
	broadcast chan *BroadcastMessage

	// システムからの全体ブロードキャストメッセージ
	systemBroadcast chan interface{}

	// ユーザーからの部屋単位ブロードキャスト
	roomBroadcast chan *RoomBroadcastMessage

	// システムからの部屋単位ブロードキャスト
	systemRoomBroadcast chan *RoomBroadcastMessage

	// 特定ユーザーへのメッセージ
	userMessage chan *UserMessage

	// 新しいクライアントの登録リクエスト
	register chan *Client

	// クライアントの登録解除リクエスト
	unregister chan *Client

	// 部屋への参加リクエスト
	joinRoom chan struct {
		UserID uuid.UUID
		RoomID uuid.UUID
	}

	// 部屋からの退出リクエスト
	leaveRoom chan struct {
		UserID uuid.UUID
		RoomID uuid.UUID
	}

	// 並行アクセス制御
	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:             make(map[uuid.UUID]*Client),
		rooms:               make(map[uuid.UUID]map[uuid.UUID]bool),
		broadcast:           make(chan *BroadcastMessage, 256),
		systemBroadcast:     make(chan interface{}, 256),
		roomBroadcast:       make(chan *RoomBroadcastMessage, 256),
		systemRoomBroadcast: make(chan *RoomBroadcastMessage, 256),
		userMessage:         make(chan *UserMessage, 256),
		register:            make(chan *Client),
		unregister:          make(chan *Client),
		joinRoom: make(chan struct {
			UserID uuid.UUID
			RoomID uuid.UUID
		}),
		leaveRoom: make(chan struct {
			UserID uuid.UUID
			RoomID uuid.UUID
		}),
	}
}

// Hubを起動
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
				// ユーザーが参加している全ての部屋から退出
				h.removeUserFromAllRooms(client.UserID)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			// ユーザーからのメッセージを全員にブロードキャスト
			h.broadcastToAll(message)

		case message := <-h.systemBroadcast:
			// システムからのメッセージを全員にブロードキャスト
			h.broadcastSystemMessage(message)

		case msg := <-h.roomBroadcast:
			// ユーザーからの部屋単位ブロードキャスト
			h.broadcastToRoom(msg.RoomID, msg.Message)

		case msg := <-h.systemRoomBroadcast:
			// システムからの部屋単位ブロードキャスト
			h.broadcastToRoom(msg.RoomID, msg.Message)

		case msg := <-h.userMessage:
			// 特定ユーザーにメッセージを送信
			h.sendToUser(msg.UserID, msg.Message)

		case req := <-h.joinRoom:
			// 部屋に参加
			h.mu.Lock()
			if h.rooms[req.RoomID] == nil {
				h.rooms[req.RoomID] = make(map[uuid.UUID]bool)
			}
			h.rooms[req.RoomID][req.UserID] = true
			h.mu.Unlock()

		case req := <-h.leaveRoom:
			// 部屋から退出
			h.mu.Lock()
			if members, ok := h.rooms[req.RoomID]; ok {
				delete(members, req.UserID)
				if len(members) == 0 {
					delete(h.rooms, req.RoomID)
				}
			}
			h.mu.Unlock()
		}
	}
}

// ユーザーを全ての部屋から削除（unregister時に呼ばれる）
func (h *Hub) removeUserFromAllRooms(userID uuid.UUID) {
	for roomID, members := range h.rooms {
		delete(members, userID)
		if len(members) == 0 {
			delete(h.rooms, roomID)
		}
	}
}

// メッセージを全クライアントに送信
func (h *Hub) broadcastToAll(msg *BroadcastMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		select {
		case client.Send <- msg.Message:
		default:
			// 送信バッファが詰まっている場合はクライアントをクローズ
			close(client.Send)
			delete(h.clients, client.UserID)
		}
	}
}

// システムメッセージを全クライアントに送信
func (h *Hub) broadcastSystemMessage(message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		select {
		case client.Send <- data:
		default:
			close(client.Send)
			delete(h.clients, client.UserID)
		}
	}
}

// 特定のユーザーにメッセージを送信（公開メソッド）
func (h *Hub) SendToUser(userID uuid.UUID, message interface{}) error {
	h.userMessage <- &UserMessage{
		UserID:  userID,
		Message: message,
	}
	return nil
}

// 特定のユーザーにメッセージを送信（内部実装）
func (h *Hub) sendToUser(userID uuid.UUID, message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()

	if !ok {
		return // クライアントが存在しない場合は無視
	}

	select {
	case client.Send <- data:
	default:
		// 送信バッファが詰まっている場合は無視
	}
}

// 部屋のメンバー全員にメッセージを送信（公開メソッド - ユーザー起点）
func (h *Hub) BroadcastToRoom(roomID uuid.UUID, message interface{}) error {
	h.roomBroadcast <- &RoomBroadcastMessage{
		RoomID:  roomID,
		Message: message,
	}
	return nil
}

// システムから部屋のメンバー全員にメッセージを送信（公開メソッド - システム起点）
func (h *Hub) SystemBroadcastToRoom(roomID uuid.UUID, message interface{}) error {
	h.systemRoomBroadcast <- &RoomBroadcastMessage{
		RoomID:  roomID,
		Message: message,
	}
	return nil
}

// 部屋のメンバー全員にメッセージを送信（内部実装）
func (h *Hub) broadcastToRoom(roomID uuid.UUID, message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	h.mu.RLock()
	members, ok := h.rooms[roomID]
	h.mu.RUnlock()

	if !ok {
		return // 部屋が存在しない
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for userID := range members {
		client, ok := h.clients[userID]
		if !ok {
			continue // クライアントが存在しない
		}

		select {
		case client.Send <- data:
		default:
			// 送信バッファが詰まっている場合は無視
		}
	}
}

// 全ユーザーにメッセージをブロードキャスト
func (h *Hub) BroadcastToAll(message interface{}) error {
	h.systemBroadcast <- message
	return nil
}

// ユーザーを部屋に参加させる（公開メソッド）
func (h *Hub) JoinRoom(userID, roomID uuid.UUID) {
	h.joinRoom <- struct {
		UserID uuid.UUID
		RoomID uuid.UUID
	}{UserID: userID, RoomID: roomID}
}

// ユーザーを部屋から退出させる（公開メソッド）
func (h *Hub) LeaveRoom(userID, roomID uuid.UUID) {
	h.leaveRoom <- struct {
		UserID uuid.UUID
		RoomID uuid.UUID
	}{UserID: userID, RoomID: roomID}
}

// 部屋のメンバー一覧を取得
func (h *Hub) GetRoomMembers(roomID uuid.UUID) []uuid.UUID {
	h.mu.RLock()
	defer h.mu.RUnlock()

	members, ok := h.rooms[roomID]
	if !ok {
		return []uuid.UUID{}
	}

	result := make([]uuid.UUID, 0, len(members))
	for userID := range members {
		result = append(result, userID)
	}
	return result
}

// 新しいクライアントを作成してHubに登録
func (h *Hub) NewClient(userID uuid.UUID, conn interface{}) *Client {
	client := &Client{
		UserID: userID,
		Hub:    h,
		Conn:   conn.(*websocket.Conn),
		Send:   make(chan []byte, 256),
	}

	h.register <- client

	return client
}
