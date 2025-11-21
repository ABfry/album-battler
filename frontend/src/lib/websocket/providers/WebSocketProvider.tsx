"use client";

import { ReactNode, useCallback, useEffect, useRef, useState } from "react";
import { WebSocketContext } from "../contexts/WebSocketContext";
import type { ConnectionStatus } from "../types";

const isWebSocketBusy = (socket: WebSocket | null) => {
  if (!socket) {
    return false;
  }
  return (
    socket.readyState === WebSocket.CONNECTING ||
    socket.readyState === WebSocket.OPEN ||
    socket.readyState === WebSocket.CLOSING
  );
};

const canCloseWebSocket = (socket: WebSocket | null) => {
  if (!socket) {
    return false;
  }
  return (
    socket.readyState === WebSocket.CONNECTING ||
    socket.readyState === WebSocket.OPEN
  );
};

type WebSocketConnection = {
  status: ConnectionStatus;
  messages: string[];
  sendMessage: (data: string) => void;
  getWebSocket: () => WebSocket | null;
  connect: () => void;
  disconnect: () => void;
};

/**
 * WebSocket接続を管理するカスタムフック（遅延接続対応）
 * @param url - WebSocketサーバのURL
 * @returns WebSocket接続の状態と操作関数
 */
function useWebSocketConnection(url: string): WebSocketConnection {
  const wsRef = useRef<WebSocket | null>(null);
  const [status, setStatus] = useState<ConnectionStatus>("idle");
  const [messages, setMessages] = useState<string[]>([]);
  const [shouldConnect, setShouldConnect] = useState(false);

  // WebSocketインスタンスを取得する関数
  const getWebSocket = useCallback(() => wsRef.current, []);

  // メッセージ送信関数
  const sendMessage = useCallback((data: string) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(data);
    } else {
      console.warn("WebSocket is not open. Cannot send message.");
    }
  }, []);

  // 接続開始
  const connect = useCallback(() => {
    if (!url) {
      console.log("WebSocket URL is empty. Skipping connection.");
      return;
    }
    if (isWebSocketBusy(wsRef.current)) {
      console.log("WebSocket is already connecting or active");
      return;
    }
    setShouldConnect(true);
    setStatus("connecting");
  }, [url]);

  // 切断
  const disconnect = useCallback(() => {
    setShouldConnect(false);
    if (wsRef.current) {
      if (canCloseWebSocket(wsRef.current)) {
        wsRef.current.close();
      }
      wsRef.current = null;
    }
    setStatus("disconnected");
  }, []);

  useEffect(() => {
    if (!shouldConnect) {
      return;
    }

    // 既に接続状態にある場合は新規接続を作成しない
    if (isWebSocketBusy(wsRef.current)) {
      return;
    }

    // WebSocket接続を作成
    const websocket = new WebSocket(url);
    wsRef.current = websocket;

    websocket.onopen = () => {
      console.log("WebSocket connected");
      setStatus("connected");
      setMessages((prev) => [...prev, "Connected to WebSocket server"]);
    };

    websocket.onmessage = (event) => {
      console.log("Message from server:", event.data);
      setMessages((prev) => [...prev, `Received: ${event.data}`]);
    };

    websocket.onerror = (error) => {
      console.error("WebSocket error:", error);
      setStatus("error");
      setMessages((prev) => [...prev, "Error occurred"]);
    };

    websocket.onclose = () => {
      console.log("WebSocket disconnected");
      if (wsRef.current === websocket) {
        wsRef.current = null;
      }
      setStatus("disconnected");
      setMessages((prev) => [...prev, "Disconnected from server"]);
    };

    // クリーンアップ
    return () => {
      if (canCloseWebSocket(websocket)) {
        websocket.close();
      }
      if (
        wsRef.current === websocket &&
        websocket.readyState === WebSocket.CLOSED
      ) {
        wsRef.current = null;
      }
    };
  }, [shouldConnect, url]);

  return {
    status,
    messages,
    sendMessage,
    getWebSocket,
    connect,
    disconnect,
  };
}

type WebSocketProviderProps = {
  url: string;
  children: ReactNode;
};

/**
 * WebSocketProvider
 * WebSocket接続を管理して，Contextを通じて子コンポーネントに状態を渡す
 * 自動接続はせず、connect()を呼ぶまで接続しない（遅延接続）
 */
export function WebSocketProvider({ url, children }: WebSocketProviderProps) {
  const connection = useWebSocketConnection(url);

  const contextValue = {
    getWebSocket: connection.getWebSocket,
    status: connection.status,
    messages: connection.messages,
    sendMessage: connection.sendMessage,
    connect: connection.connect,
    disconnect: connection.disconnect,
  };

  return (
    <WebSocketContext.Provider value={contextValue}>
      {children}
    </WebSocketContext.Provider>
  );
}
