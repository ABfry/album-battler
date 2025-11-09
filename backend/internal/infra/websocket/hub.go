package websocket

import (
	"encoding/json"
	"log"
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
				client.Close()
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
	// クライアントリストをコピー（スナップショット）
	clients := make([]*Client, 0, len(h.clients))
	for _, client := range h.clients {
		if !client.IsClosed() {
			clients = append(clients, client)
		}
	}
	h.mu.RUnlock()

	// ロック外で送信処理
	var toDelete []uuid.UUID
	for _, client := range clients {
		if !h.trySendToClient(client, msg.Message) {
			toDelete = append(toDelete, client.UserID)
		}
	}

	// 削除が必要な場合
	if len(toDelete) > 0 {
		h.removeClients(toDelete)
	}
}

// システムメッセージを全クライアントに送信
func (h *Hub) broadcastSystemMessage(message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	h.mu.RLock()
	// 削除対象を記録
	var toUnregister []*Client
	for _, client := range h.clients {
		select {
		case client.Send <- data:
		default:
			toUnregister = append(toUnregister, client)
		}
	}
	h.mu.RUnlock()

	// unregisterチャネル経由で削除
	for _, client := range toUnregister {
		select {
		case h.unregister <- client:
		default:
			log.Printf("Failed to unregister client %s", client.UserID)
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
	defer h.mu.RUnlock()

	members, ok := h.rooms[roomID]
	if !ok {
		return
	}

	for userID := range members {
		if client, ok := h.clients[userID]; ok {
			select {
			case client.Send <- data:
			default:
				// バッファフル時は無視
			}
		}
	}
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

// 送信メソッド
func (h *Hub) trySendToClient(client *Client, msg []byte) bool {
	// 既にcloseされているかチェック
	if client.IsClosed() {
		return false
	}

	select {
	case client.Send <- msg:
		return true
	default:
		// バッファフルの場合
		client.Close() // sync.Onceで保護
		return false
	}
}

// クライアント削除
func (h *Hub) removeClients(userIDs []uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, userID := range userIDs {
		if client, ok := h.clients[userID]; ok {
			client.Close()
			delete(h.clients, userID)

			// roomsから削除
			for roomID, members := range h.rooms {
				if members[userID] {
					delete(members, userID)
					if len(members) == 0 {
						delete(h.rooms, roomID)
					}
				}
			}

			log.Printf("Client %s disconnected and cleaned up", userID)
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
