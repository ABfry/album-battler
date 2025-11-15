// src/lib/websocket/hooks/useWebSocketEvents.ts
import { useCallback } from "react";
import { useWebSocket } from "../contexts/WebSocketContext";
import type { EventMap, WebSocketMessage } from "../types";

export type EventHandler<T = unknown> = (
  payload: T,
  message: WebSocketMessage<T>
) => void;

export function useWebSocketEvents() {
  const { getWebSocket } = useWebSocket();

  const subscribe = useCallback(
    <K extends keyof EventMap>(
      eventType: K,
      handler: EventHandler<EventMap[K]>
    ) => {
    const ws = getWebSocket();
    if (!ws) {
      console.warn("WebSocket not connected");
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
    [getWebSocket]
  );

  return { subscribe };
}
