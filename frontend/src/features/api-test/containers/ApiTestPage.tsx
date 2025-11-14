"use client";

import { useState, useEffect } from "react";
import { useRoom } from "@/src/hooks/useRoom";
import { useRoomInfo } from "@/src/hooks/useRoomInfo";
import { useWebSocket } from "@/src/lib/websocket/contexts/WebSocketContext";
import { useWebSocketEvents } from "@/src/lib/websocket/hooks/useWebSocketEvents";
import { ApiTestView } from "../components/ApiTestView";
import type { CreateRoomResponse } from "@/src/lib/api/types";
import type {
  PlayerJoinRoomPayload,
  PlayerLeaveRoomPayload,
  StartGamePayload,
} from "@/src/lib/websocket/types";

/**
 * API テストページ Container
 * ロジック・状態管理を担当
 */
export function ApiTestPage() {
  // 部屋作成・参加・退出・ゲーム開始の操作
  const { createRoom, joinRoom, leaveRoom, startGame, loading, error } =
    useRoom();

  // WebSocket接続
  const { connect, disconnect } = useWebSocket();
  const { subscribe } = useWebSocketEvents();

  // 作成された部屋の情報を保持
  const [createdRoom, setCreatedRoom] = useState<CreateRoomResponse | null>(
    null
  );

  // 部屋情報の取得（作成後に自動取得）
  const {
    room,
    loading: roomLoading,
    error: roomError,
    refetch,
  } = useRoomInfo(createdRoom?.room_id || null);

  // 各操作の成功状態
  const [joinSuccess, setJoinSuccess] = useState(false);
  const [leaveSuccess, setLeaveSuccess] = useState(false);
  const [startSuccess, setStartSuccess] = useState(false);

  // WebSocket接続の初期化
  useEffect(() => {
    connect();
    return () => disconnect();
  }, [connect, disconnect]);

  // WebSocketイベントの購読
  useEffect(() => {
    const unsubscribeJoin = subscribe(
      "player_join_room",
      (payload: PlayerJoinRoomPayload) => {
        console.log("Player joined room:", payload.room_id);
        refetch();
      }
    );

    const unsubscribeLeave = subscribe(
      "player_leave_room",
      (payload: PlayerLeaveRoomPayload) => {
        console.log("Player left room:", payload.room_id);
        refetch();
      }
    );

    const unsubscribeStart = subscribe(
      "start_game",
      (payload: StartGamePayload) => {
        console.log("Game started:", payload.room_id);
        refetch();
      }
    );

    return () => {
      unsubscribeJoin();
      unsubscribeLeave();
      unsubscribeStart();
    };
  }, [subscribe, refetch]);

  // 1. 部屋作成
  const handleCreateRoom = async (userId: string) => {
    const result = await createRoom(userId);
    if (result) {
      setCreatedRoom(result);
      setJoinSuccess(false);
      setStartSuccess(false);
    }
  };

  // 2. 部屋参加
  const handleJoinRoom = async (userId: string, roomNumber: number) => {
    const success = await joinRoom(userId, roomNumber);
    setJoinSuccess(success);
  };

  // 3. 部屋退出
  const handleLeaveRoom = async (roomId: string, userId: string) => {
    const success = await leaveRoom(roomId, userId);
    setLeaveSuccess(success);
    if (success) {
      setTimeout(() => refetch(), 500);
    }
  };

  // 4. ゲーム開始
  const handleStartGame = async () => {
    if (!createdRoom) return;
    const success = await startGame(
      createdRoom.room_id,
      "550e8400-e29b-41d4-a716-446655440001"
    );
    setStartSuccess(success);
  };

  return (
    <ApiTestView
      createdRoom={createdRoom}
      room={room}
      loading={loading}
      roomLoading={roomLoading}
      error={error}
      roomError={roomError}
      joinSuccess={joinSuccess}
      leaveSuccess={leaveSuccess}
      startSuccess={startSuccess}
      onCreateRoom={handleCreateRoom}
      onJoinRoom={handleJoinRoom}
      onLeaveRoom={handleLeaveRoom}
      onStartGame={handleStartGame}
      onRefetch={refetch}
    />
  );
}
