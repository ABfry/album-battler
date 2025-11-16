"use client";

import { useState } from "react";
import { battleApi } from "@/src/lib/api/battleApi";
import type { CreateBattleResponse } from "@/src/lib/api/types";

type UseBattleResult = {
  createBattle: (roomId: string) => Promise<CreateBattleResponse | null>;
  sendImage: (
    battleId: string,
    userId: string,
    imageBase64: string
  ) => Promise<boolean>;
  loading: boolean;
  error: string | null;
};

/**
 * バトル操作用のカスタムフック
 * バトル作成、画像送信の処理を管理
 */
export function useBattle(): UseBattleResult {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const createBattle = async (roomId: string) => {
    setLoading(true);
    setError(null);
    try {
      const result = await battleApi.createBattle(roomId);
      return result;
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to create battle";
      setError(message);
      return null;
    } finally {
      setLoading(false);
    }
  };

  const sendImage = async (
    battleId: string,
    userId: string,
    imageBase64: string
  ) => {
    setLoading(true);
    setError(null);
    try {
      await battleApi.sendImage(battleId, userId, imageBase64);
      return true;
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to send image";
      setError(message);
      return false;
    } finally {
      setLoading(false);
    }
  };

  return {
    createBattle,
    sendImage,
    loading,
    error,
  };
}
