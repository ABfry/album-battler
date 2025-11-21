// リクエスト型
export type CreateRoomRequest = {
  user_id: string;
};

export type JoinRoomRequest = {
  user_id: string;
  room_number: number;
};

export type StartGameRequest = {
  user_id: string;
};

export type LeaveRoomRequest = {
  user_id: string;
};

// レスポンス型
export type CreateRoomResponse = {
  room_id: string;
  room_number: number;
};

export type JoinRoomResponse = {
  message: string;
};

export type UserInfo = {
  id: string;
  name: string;
  icon_url: string;
};

export type RoomInfoResponse = {
  room_number: number;
  users: UserInfo[];
  is_expired: boolean;
  room_status: RoomStatusString;
  host_user_id: string | null;
};

export type StartGameResponse = {
  result: string;
};

export type LeaveRoomResponse = {
  message: string;
};

export type GetBattleIDResponse = {
  BattleID: string;
};

// Battle リクエスト型
export type CreateBattleRequest = {
  room_id: string;
};

export type SendImageRequest = {
  user_id: string;
  image_base64: string;
};

// Battle レスポンス型
export type CreateBattleResponse = {
  BattleID: string;
};

export type GetBattleResponse = {
  Battle: {
    ID: string;
    RoomID: string;
    StartedAt: string;
    Theme: string;
    UserIDs: string[];
  };
};

export type GetImageResponse = {
  Images: Array<{
    UserID: string;
    ImageURL: string;
  }>;
};

export type SendImageResponse = {
  result: string;
};

export type UserResult = {
  user_id: string;
  ai_score: number;
  user_score: number;
  final_score: number;
  rank: number;
  image_url: string;
  ai_explanation: string;
};

export type GetBattleResultResponse = {
  battle_id: string;
  winner_user_id: string;
  results: UserResult[];
};

// ドメイン型
export type RoomStatusString =
  | "waiting"
  | "full"
  | "battling"
  | "result"
  | "closed";

export enum RoomStatus {
  WaitJoin = "waiting",
  FullyJoined = "full",
  InBattle = "battling",
  Result = "result",
  Closed = "closed",
}

export type Room = {
  id: string;
  roomNumber: number;
  hostUserId: string | null;
  users: UserInfo[];
  status: RoomStatus;
  isExpired: boolean;
};

export type Battle = {
  id: string;
  roomId: string;
  startedAt: string;
  theme: string;
  userIds: string[];
};

export type BattleImage = {
  userId: string;
  imageUrl: string;
};

// WebSocketイベント型
export type UserJoinedRoomEvent = {
  type: "user_joined_room";
  room_id: string;
  room_number: number;
  user_id: string;
  is_host: boolean;
  occurred_at: string;
};

export type UserLeftRoomEvent = {
  type: "user_left_room";
  room_id: string;
  room_number: number;
  user_id: string;
  was_host: boolean;
  new_host_id: string | null;
  room_dissolved: boolean;
  occurred_at: string;
};

export type GameStartedEvent = {
  type: "game_started";
  room_id: string;
  room_number: number;
  occurred_at: string;
};

export type WebSocketEvent =
  | UserJoinedRoomEvent
  | UserLeftRoomEvent
  | GameStartedEvent;
