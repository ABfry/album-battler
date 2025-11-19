import { useCallback } from "react";
import { useWebSocket } from "../contexts/WebSocketContext";

type SendClapParams = {
  userId: string;
  targetUserId: string;
  battleId: string;
  count: number;
};

/**
 * WebSocket経由で拍手を送信するhook
 * 既存のWebSocket接続を共有して使用
 */
export function useWebSocketClap() {
  const { sendMessage, status } = useWebSocket();

  /**
   * 拍手を送信
   * @param params - userId, targetUserId, battleId, count (1以上)
   * @throws Error - countが1未満、またはWebSocket未接続の場合
   */
  const sendClap = useCallback(
    ({ userId, targetUserId, battleId, count }: SendClapParams) => {
      // バリデーション
      if (count < 1) {
        throw new Error("Clap count must be at least 1");
      }

      if (status !== "connected") {
        console.error(
          `[useWebSocketClap] Cannot send clap: WebSocket status is ${status}`
        );
        throw new Error("WebSocket is not connected");
      }

      const message = {
        type: "clap_send",
        payload: {
          user_id: userId,
          target_user_id: targetUserId,
          battle_id: battleId,
          count,
        },
      };

      console.log("[useWebSocketClap] Sending clap:", message);
      sendMessage(JSON.stringify(message));
    },
    [sendMessage, status]
  );

  return {
    sendClap,
    isConnected: status === "connected",
    connectionStatus: status,
  };
}
