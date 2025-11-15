export type ConnectionStatus =
  | "idle"
  | "connecting"
  | "connected"
  | "disconnected"
  | "error";

export type WebSocketMessage<T = unknown> = {
  type: string;
  payload: T;
  timestamp: string;
};

// バックエンドから送信されるイベント型
export type PlayerJoinRoomPayload = {
  room_id: string;
};

export type PlayerLeaveRoomPayload = {
  room_id: string;
};

export type StartGamePayload = {
  room_id: string;
};

// イベントマップ（型安全）
export type EventMap = {
  player_join_room: PlayerJoinRoomPayload;
  player_leave_room: PlayerLeaveRoomPayload;
  start_game: StartGamePayload;
};
