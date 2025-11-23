// src/lib/websocket/hooks/useWebSocketEvents.ts
import { useCallback } from "react";
import { useWebSocket } from "../contexts/WebSocketContext";
import type { EventMap, WebSocketMessage } from "../types";

export type EventHandler<T = unknown> = (
  payload: T,
  message: WebSocketMessage<T>
) => void;

export function useWebSocketEvents() {
  const { getWebSocket, status, connect } = useWebSocket();

  const subscribe = useCallback(
    <K extends keyof EventMap>(
      eventType: K,
      handler: EventHandler<EventMap[K]>
    ) => {
      const ws = getWebSocket();
      if (!ws) {
        console.warn("WebSocket not connected, attempting to reconnect...");
        // 切断状態またはアイドル状態の場合は再接続を試みる
        if (status === "disconnected" || status === "idle") {
          // console.log("Triggering reconnection from useWebSocketEvents");
          connect();
        }
        return () => {};
      }

      const listener = (event: MessageEvent) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data);

          // イベントタイプが一致したら処理
          if (message.type === eventType) {
            const typedMessage = message as WebSocketMessage<EventMap[K]>;
            handler(message.payload as EventMap[K], typedMessage);
          }
        } catch (error) {
          console.error("Failed to parse WebSocket message", error);
        }
      };

      ws.addEventListener("message", listener);

      // クリーンアップ関数を返す
      return () => {
        ws.removeEventListener("message", listener);
      };
    },
    [getWebSocket, status, connect]
  );

  return { subscribe };
}
