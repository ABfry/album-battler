"use client";

import { useState, useEffect, useCallback } from "react";
import { battleApi } from "@/src/lib/api/battleApi";
import type {
  Battle,
  BattleImage,
  GetBattleResponse,
  GetImageResponse,
} from "@/src/lib/api/types";

type UseBattleInfoResult = {
  battle: Battle | null;
  images: BattleImage[];
  loading: boolean;
  error: string | null;
  refetch: () => Promise<void>;
  refetchImages: () => Promise<void>;
};

/**
 * GetBattleResponseをフロントエンド用のBattle型に変換
 */
function mapToBattle(data: GetBattleResponse): Battle {
  return {
    id: data.Battle.ID,
    roomId: data.Battle.RoomID,
    startedAt: data.Battle.StartedAt,
    theme: data.Battle.Theme,
    userIds: data.Battle.UserIDs,
    battleTimeLimitSeconds: data.Battle.BattleTimeLimitSeconds,
  };
}

/**
 * GetImageResponseをフロントエンド用のBattleImage配列に変換
 */
function mapToImages(data: GetImageResponse): BattleImage[] {
  return data.Images.map((img) => ({
    userId: img.UserID,
    imageUrl: img.ImageURL,
  }));
}

/**
 * バトル情報取得用のカスタムフック
 * バトル情報と画像一覧を取得し、リアルタイムで更新可能
 */
export function useBattleInfo(battleId: string | null): UseBattleInfoResult {
  const [battle, setBattle] = useState<Battle | null>(null);
  const [images, setImages] = useState<BattleImage[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchBattleInfo = useCallback(async () => {
    if (!battleId) {
      setBattle(null);
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const data = await battleApi.getBattle(battleId);
      setBattle(mapToBattle(data));
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to fetch battle info";
      setError(message);
      setBattle(null);
    } finally {
      setLoading(false);
    }
  }, [battleId]);

  const fetchImages = useCallback(async () => {
    if (!battleId) {
      setImages([]);
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const data = await battleApi.getImages(battleId);
      setImages(mapToImages(data));
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to fetch images";
      setError(message);
      setImages([]);
    } finally {
      setLoading(false);
    }
  }, [battleId]);

  useEffect(() => {
    fetchBattleInfo();
    fetchImages(); // 画像一覧も初回に取得
  }, [fetchBattleInfo, fetchImages]);

  return {
    battle,
    images,
    loading,
    error,
    refetch: fetchBattleInfo,
    refetchImages: fetchImages,
  };
}
