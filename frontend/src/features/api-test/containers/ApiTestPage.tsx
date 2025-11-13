"use client";

import { useState } from "react";
import { useRoom } from "@/src/hooks/useRoom";
import { useRoomInfo } from "@/src/hooks/useRoomInfo";
import { ApiTestView } from "../components/ApiTestView";
import type { CreateRoomResponse } from "@/src/lib/api/types";

/**
 * API テストページ Container
 * ロジック・状態管理を担当
 */
export function ApiTestPage() {
  // 部屋作成・参加・ゲーム開始の操作
  const { createRoom, joinRoom, startGame, loading, error } = useRoom();

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
  const [startSuccess, setStartSuccess] = useState(false);

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
    if (success) {
      setTimeout(() => refetch(), 500);
    }
  };

  // 3. ゲーム開始
  const handleStartGame = async () => {
    if (!createdRoom) return;
    const success = await startGame(
      createdRoom.room_id,
      "550e8400-e29b-41d4-a716-446655440001"
    );
    setStartSuccess(success);
    if (success) {
      setTimeout(() => refetch(), 500);
    }
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
      startSuccess={startSuccess}
      onCreateRoom={handleCreateRoom}
      onJoinRoom={handleJoinRoom}
      onStartGame={handleStartGame}
      onRefetch={refetch}
    />
  );
}
