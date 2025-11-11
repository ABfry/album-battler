"use client";

import { ReactNode } from "react";
import { WebSocketContext } from "../contexts/WebSocketContext";
import { useWebSocketConnection } from "../hooks/useWebSocketConnection";

type WebSocketProviderProps = {
  url: string;
  children: ReactNode;
};

/**
 * WebSocketProvider
 * WebSocket接続を管理して，Contextを通じて子コンポーネントに状態を渡す
 */
export function WebSocketProvider({ url, children }: WebSocketProviderProps) {
  const connection = useWebSocketConnection(url);

  const contextValue = {
    getWebSocket: connection.getWebSocket,
    status: connection.status,
    messages: connection.messages,
    sendMessage: connection.sendMessage,
  };

  return (
    <WebSocketContext.Provider value={contextValue}>
      {children}
    </WebSocketContext.Provider>
  );
}
