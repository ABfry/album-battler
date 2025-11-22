"use client";

import { useEffect, useState } from "react";
import { useWebSocket } from "@/src/lib/websocket/contexts/WebSocketContext";
import { roomApi } from "@/src/lib/api/roomApi";
import { battleApi } from "@/src/lib/api/battleApi";
import { getUserIdClient } from "@/src/lib/auth/getUserIdClient";
import type { GetBattleStateResponse } from "@/src/lib/api/types";

type GamePhase = "selecting" | "clap_time" | "result" | "finished";

type GameStateRestoreResult = {
  isRestoring: boolean;
  error: string | null;
  restoredState: GetBattleStateResponse | null;
};

/**
 * ゲーム状態復元用カスタムフック
 * WebSocket再接続時に部屋への再参加とゲーム状態の復元を行う
 */
export function useGameStateRestore(): GameStateRestoreResult {
  const { status, isReconnecting, currentRoomId, currentBattleId } =
    useWebSocket();
  const [isRestoring, setIsRestoring] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [restoredState, setRestoredState] =
    useState<GetBattleStateResponse | null>(null);

  useEffect(() => {
    // 再接続時のみ処理を実行
    if (
      status === "connected" &&
      isReconnecting &&
      currentRoomId &&
      currentBattleId
    ) {
      const restoreGameState = async () => {
        setIsRestoring(true);
        setError(null);

        try {
          // ユーザーIDを取得
          const userId = getUserIdClient();
          if (!userId) {
            throw new Error("User ID not found");
          }

          // 1. 部屋に再参加（WebSocket Hubのroomsに再登録）
          console.log("Rejoining room:", currentRoomId);
          await roomApi.rejoinRoom(userId, currentRoomId);

          // 2. バトルの状態を取得
          console.log("Fetching battle state:", currentBattleId);
          const battleState = await battleApi.getBattleState(currentBattleId);

          console.log("Battle state restored:", battleState);
          setRestoredState(battleState);
        } catch (err) {
          const message =
            err instanceof Error
              ? err.message
              : "Failed to restore game state";
          console.error("Game state restoration error:", message);
          setError(message);
        } finally {
          setIsRestoring(false);
        }
      };

      restoreGameState();
    }
  }, [status, isReconnecting, currentRoomId, currentBattleId]);

  return {
    isRestoring,
    error,
    restoredState,
  };
}
