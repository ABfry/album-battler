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

export type ImageSendPayload = {
  room_id: string;
  battle_id: string;
  user_id: string;
  image_url: string;
};

export type StartClapTimePayload = Record<string, never>; // バックエンドは空のペイロードを送信

export type ClapSendPayload = Record<string, never>; // バックエンドは空のペイロードを送信（誰が何回拍手したかは含まれない）

export type StartResultPhasePayload = {
  room_id: string;
  battle_id: string;
  winner_user_id: string;
};

export type RoomSettingsUpdatedPayload = {
  room_id: string;
  room_number: number;
  battle_time_limit_seconds: number;
};

// イベントマップ（型安全）
export type EventMap = {
  player_join_room: PlayerJoinRoomPayload;
  player_leave_room: PlayerLeaveRoomPayload;
  start_game: StartGamePayload;
  image_send: ImageSendPayload;
  start_clap_time: StartClapTimePayload;
  clap_send: ClapSendPayload;
  start_result_phase: StartResultPhasePayload;
  room_settings_updated: RoomSettingsUpdatedPayload;
};
