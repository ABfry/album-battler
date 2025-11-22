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
  setUrl: (url: string) => void;
  currentUrl: string;
  isReconnecting: boolean;
};

/**
 * WebSocket接続を管理するカスタムフック（遅延接続対応）
 * @param initialUrl - WebSocketサーバの初期URL
 * @returns WebSocket接続の状態と操作関数
 */
function useWebSocketConnection(initialUrl: string): WebSocketConnection {
  const wsRef = useRef<WebSocket | null>(null);
  const [status, setStatus] = useState<ConnectionStatus>("idle");
  const [messages, setMessages] = useState<string[]>([]);
  const [shouldConnect, setShouldConnect] = useState(false);
  const [currentUrl, setCurrentUrl] = useState(initialUrl);
  const [isReconnecting, setIsReconnecting] = useState(false);

  // 再接続関連の状態
  const reconnectAttemptsRef = useRef(0);
  const maxReconnectAttempts = 5;
  const isManualDisconnectRef = useRef(false);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  // WebSocketインスタンスを取得する関数
  const getWebSocket = useCallback(() => wsRef.current, []);

  // 再接続タイマーをクリア
  const clearReconnectTimeout = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }
  }, []);

  // メッセージ送信関数
  const sendMessage = useCallback((data: string) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(data);
    } else {
      console.warn("WebSocket is not open. Cannot send message.");
    }
  }, []);

  // URL を動的に設定
  const setUrl = useCallback(
    (newUrl: string) => {
      if (newUrl === currentUrl) {
        return;
      }
      // 既存の接続があれば切断
      if (canCloseWebSocket(wsRef.current)) {
        wsRef.current?.close();
        wsRef.current = null;
      }
      clearReconnectTimeout();
      reconnectAttemptsRef.current = 0;
      setShouldConnect(false);
      setCurrentUrl(newUrl);
      console.log("WebSocket URL updated:", newUrl);
    },
    [currentUrl, clearReconnectTimeout]
  );

  // 接続開始
  const connect = useCallback(() => {
    if (!currentUrl) {
      console.log("WebSocket URL is empty. Skipping connection.");
      return;
    }
    if (isWebSocketBusy(wsRef.current)) {
      console.log("WebSocket is already connecting or active");
      return;
    }
    isManualDisconnectRef.current = false;
    clearReconnectTimeout();
    setShouldConnect(true);
    setStatus("connecting");
  }, [currentUrl, clearReconnectTimeout]);

  // 切断
  const disconnect = useCallback(() => {
    isManualDisconnectRef.current = true;
    clearReconnectTimeout();
    reconnectAttemptsRef.current = 0;
    setShouldConnect(false);
    if (wsRef.current) {
      if (canCloseWebSocket(wsRef.current)) {
        wsRef.current.close();
      }
      wsRef.current = null;
    }
    setStatus("disconnected");
  }, [clearReconnectTimeout]);

  useEffect(() => {
    if (!shouldConnect) {
      return;
    }

    // 既に接続状態にある場合は新規接続を作成しない
    if (isWebSocketBusy(wsRef.current)) {
      return;
    }

    // WebSocket接続を作成
    const websocket = new WebSocket(currentUrl);
    wsRef.current = websocket;

    websocket.onopen = () => {
      console.log("WebSocket connected");
      const wasReconnecting = reconnectAttemptsRef.current > 0;
      reconnectAttemptsRef.current = 0;
      setStatus("connected");
      setIsReconnecting(wasReconnecting);
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

      // 手動切断でない場合のみ再接続を試みる
      if (!isManualDisconnectRef.current && shouldConnect) {
        if (reconnectAttemptsRef.current < maxReconnectAttempts) {
          const delay = Math.min(
            1000 * Math.pow(2, reconnectAttemptsRef.current),
            32000
          );
          console.log(
            `WebSocket will attempt to reconnect in ${delay}ms (attempt ${reconnectAttemptsRef.current + 1}/${maxReconnectAttempts})`
          );

          reconnectTimeoutRef.current = setTimeout(() => {
            reconnectAttemptsRef.current += 1;
            console.log(
              `Attempting to reconnect... (${reconnectAttemptsRef.current}/${maxReconnectAttempts})`
            );
            setShouldConnect(false); // 強制的に再接続をトリガー
            setShouldConnect(true);
          }, delay);
        } else {
          console.log("Max reconnect attempts reached. Giving up.");
          setStatus("error");
          setMessages((prev) => [
            ...prev,
            "Failed to reconnect after multiple attempts",
          ]);
        }
      }
    };

    // クリーンアップ
    return () => {
      clearReconnectTimeout();
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
  }, [shouldConnect, currentUrl, clearReconnectTimeout]);

  return {
    status,
    messages,
    sendMessage,
    getWebSocket,
    connect,
    disconnect,
    setUrl,
    currentUrl,
    isReconnecting,
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

  // 再接続時の状態復元用
  const [currentRoomId, setCurrentRoomId] = useState<string | null>(null);
  const [currentBattleId, setCurrentBattleId] = useState<string | null>(null);

  const contextValue = {
    getWebSocket: connection.getWebSocket,
    status: connection.status,
    messages: connection.messages,
    sendMessage: connection.sendMessage,
    connect: connection.connect,
    disconnect: connection.disconnect,
    setUrl: connection.setUrl,
    isReconnecting: connection.isReconnecting,
    currentRoomId,
    currentBattleId,
    setCurrentRoomId,
    setCurrentBattleId,
  };

  return (
    <WebSocketContext.Provider value={contextValue}>
      {children}
    </WebSocketContext.Provider>
  );
}
