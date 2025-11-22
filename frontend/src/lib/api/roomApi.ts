import { fetchApi } from "./client";
import type {
  CreateRoomRequest,
  CreateRoomResponse,
  GetBattleIDResponse,
  JoinRoomRequest,
  JoinRoomResponse,
  LeaveRoomRequest,
  LeaveRoomResponse,
  RejoinRoomRequest,
  RejoinRoomResponse,
  RoomInfoResponse,
  StartGameRequest,
  StartGameResponse,
} from "./types";

/**
 * 部屋関連のAPI
 */
export const roomApi = {
  /**
   * 部屋を作成
   * POST /room
   */
  createRoom: async (userId: string): Promise<CreateRoomResponse> => {
    return fetchApi<CreateRoomResponse>("/room", {
      method: "POST",
      body: JSON.stringify({ user_id: userId } satisfies CreateRoomRequest),
    });
  },

  /**
   * 部屋に参加
   * POST /room/join
   */
  joinRoom: async (
    userId: string,
    roomNumber: number
  ): Promise<JoinRoomResponse> => {
    return fetchApi<JoinRoomResponse>("/room/join", {
      method: "POST",
      body: JSON.stringify({
        user_id: userId,
        room_number: roomNumber,
      } satisfies JoinRoomRequest),
    });
  },

  /**
   * 部屋情報を取得
   * GET /room/{id}
   */
  getRoomInfo: async (roomId: string): Promise<RoomInfoResponse> => {
    return fetchApi<RoomInfoResponse>(`/room/${roomId}`);
  },

  /**
   * ゲームを開始
   * POST /room/{id}/start
   */
  startGame: async (
    roomId: string,
    userId: string
  ): Promise<StartGameResponse> => {
    return fetchApi<StartGameResponse>(`/room/${roomId}/start`, {
      method: "POST",
      body: JSON.stringify({ user_id: userId } satisfies StartGameRequest),
    });
  },

  /**
   * 部屋から退出
   * POST /room/{id}/leave
   */
  leaveRoom: async (
    roomId: string,
    userId: string
  ): Promise<LeaveRoomResponse> => {
    return fetchApi<LeaveRoomResponse>(`/room/${roomId}/leave`, {
      method: "POST",
      body: JSON.stringify({ user_id: userId } satisfies LeaveRoomRequest),
    });
  },

  /**
   * 部屋のバトルIDを取得
   * GET /room/{id}/battle-id
   */
  getBattleID: async (roomId: string): Promise<GetBattleIDResponse> => {
    return fetchApi<GetBattleIDResponse>(`/room/${roomId}/battle-id`);
  },

  /**
   * 部屋に再参加（WebSocket再接続時）
   * POST /room/rejoin
   */
  rejoinRoom: async (
    userId: string,
    roomId: string
  ): Promise<RejoinRoomResponse> => {
    return fetchApi<RejoinRoomResponse>("/room/rejoin", {
      method: "POST",
      body: JSON.stringify({
        user_id: userId,
        room_id: roomId,
      } satisfies RejoinRoomRequest),
    });
  },
};
