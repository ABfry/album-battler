import { createContext, useContext } from "react";
import type { ConnectionStatus } from "../hooks/useWebSocketConnection";

export type WebSocketContextType = {
  getWebSocket: () => WebSocket | null;
  status: ConnectionStatus;
  messages: string[];
  sendMessage: (data: string) => void;
};

export const WebSocketContext = createContext<WebSocketContextType | null>(
  null
);

/**
 * WebSocketContextを利用するためのカスタムフック
 * @throws Context外で使用したらエラー
 */
export function useWebSocket(): WebSocketContextType {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error("useWebSocket must be used within WebSocketProvider");
  }
  return context;
}
