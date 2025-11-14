"use client";

import { useState, useEffect, useCallback } from "react";
import { roomApi } from "@/src/lib/api/roomApi";
import { RoomStatus } from "@/src/lib/api/types";
import type { Room, RoomInfoResponse } from "@/src/lib/api/types";

type UseRoomInfoResult = {
  room: Room | null;
  loading: boolean;
  error: string | null;
  refetch: () => Promise<void>;
};

/**
 * 文字列をRoomStatus enumに安全に変換
 */
const statusMap: Record<string, RoomStatus> = {
  waiting: RoomStatus.WaitJoin,
  full: RoomStatus.FullyJoined,
  battling: RoomStatus.InBattle,
  result: RoomStatus.Result,
  closed: RoomStatus.Closed,
};

/**
 * RoomInfoResponseをフロントエンド用のRoom型に変換
 */
function mapToRoom(roomId: string, data: RoomInfoResponse): Room {
  return {
    id: roomId,
    roomNumber: data.room_number,
    hostUserId: data.host_user_id,
    users: data.users,
    status: statusMap[data.room_status] ?? RoomStatus.WaitJoin,
    isExpired: data.is_expired,
  };
}

/**
 * 部屋情報取得用のカスタムフック
 * 部屋情報を取得し、リアルタイムで更新可能
 */
export function useRoomInfo(roomId: string | null): UseRoomInfoResult {
  const [room, setRoom] = useState<Room | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchRoomInfo = useCallback(async () => {
    if (!roomId) {
      setRoom(null);
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const data = await roomApi.getRoomInfo(roomId);
      setRoom(mapToRoom(roomId, data));
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to fetch room info";
      setError(message);
      setRoom(null);
    } finally {
      setLoading(false);
    }
  }, [roomId]);

  useEffect(() => {
    fetchRoomInfo();
  }, [fetchRoomInfo]);

  return { room, loading, error, refetch: fetchRoomInfo };
}
