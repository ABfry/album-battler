"use client";

import { useState } from "react";
import { roomApi } from "@/src/lib/api/roomApi";
import type {
  CreateRoomResponse,
  UpdateRoomSettingsRequest,
} from "@/src/lib/api/types";

type UseRoomResult = {
  createRoom: (userId: string) => Promise<CreateRoomResponse | null>;
  joinRoom: (userId: string, roomNumber: number) => Promise<string | null>;
  leaveRoom: (roomId: string, userId: string) => Promise<boolean>;
  startGame: (roomId: string, userId: string) => Promise<boolean>;
  getBattleID: (roomId: string) => Promise<string | null>;
  updateRoomSettings: (
    roomId: string,
    userId: string,
    settings: Partial<Omit<UpdateRoomSettingsRequest, "user_id">>
  ) => Promise<boolean>;
  loading: boolean;
  error: string | null;
};

/**
 * 部屋操作用のカスタムフック
 * 部屋作成、参加、ゲーム開始の処理を管理
 */
export function useRoom(): UseRoomResult {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const createRoom = async (userId: string) => {
    setLoading(true);
    setError(null);
    try {
      const result = await roomApi.createRoom(userId);
      return result;
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to create room";
      setError(message);
      return null;
    } finally {
      setLoading(false);
    }
  };

  const joinRoom = async (userId: string, roomNumber: number) => {
    setLoading(true);
    setError(null);
    try {
      const result = await roomApi.joinRoom(userId, roomNumber);
      return result.room_id;
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to join room";
      setError(message);
      return null;
    } finally {
      setLoading(false);
    }
  };

  const leaveRoom = async (roomId: string, userId: string) => {
    setLoading(true);
    setError(null);
    try {
      await roomApi.leaveRoom(roomId, userId);
      return true;
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to leave room";
      setError(message);
      return false;
    } finally {
      setLoading(false);
    }
  };

  const startGame = async (roomId: string, userId: string) => {
    setLoading(true);
    setError(null);
    try {
      await roomApi.startGame(roomId, userId);
      return true;
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to start game";
      setError(message);
      return false;
    } finally {
      setLoading(false);
    }
  };

  const getBattleID = async (roomId: string) => {
    setLoading(true);
    setError(null);
    try {
      const result = await roomApi.getBattleID(roomId);
      return result.BattleID;
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to get battle ID";
      setError(message);
      return null;
    } finally {
      setLoading(false);
    }
  };

  const updateRoomSettings = async (
    roomId: string,
    userId: string,
    settings: Partial<Omit<UpdateRoomSettingsRequest, "user_id">>
  ) => {
    setLoading(true);
    setError(null);
    try {
      await roomApi.updateRoomSettings(roomId, userId, settings);
      return true;
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to update room settings";
      setError(message);
      return false;
    } finally {
      setLoading(false);
    }
  };

  return {
    createRoom,
    joinRoom,
    leaveRoom,
    startGame,
    getBattleID,
    updateRoomSettings,
    loading,
    error,
  };
}
