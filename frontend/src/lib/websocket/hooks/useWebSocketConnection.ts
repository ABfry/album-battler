import { useCallback, useEffect, useRef, useState } from "react";

export type ConnectionStatus =
  | "idle"
  | "connecting"
  | "connected"
  | "disconnected"
  | "error";

export type WebSocketConnection = {
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
export function useWebSocketConnection(url: string): WebSocketConnection {
  const wsRef = useRef<WebSocket | null>(null);
  const [status, setStatus] = useState<ConnectionStatus>("idle");
  const [messages, setMessages] = useState<string[]>([]);
  const shouldConnectRef = useRef(false);

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
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      console.log("WebSocket already connected");
      return;
    }
    shouldConnectRef.current = true;
    setStatus("connecting");
  }, []);

  // 切断
  const disconnect = useCallback(() => {
    shouldConnectRef.current = false;
    if (wsRef.current) {
      // 接続中または接続済みの場合のみclose()を呼ぶ
      const currentState = wsRef.current.readyState;
      if (
        currentState === WebSocket.CONNECTING ||
        currentState === WebSocket.OPEN
      ) {
        wsRef.current.close();
      }
      wsRef.current = null;
    }
    setStatus("disconnected");
  }, []);

  useEffect(() => {
    if (!shouldConnectRef.current) {
      return;
    }

    // 既に接続中または接続済みなら何もしない（重複接続を防ぐ）
    const currentState = wsRef.current?.readyState;
    if (
      currentState === WebSocket.CONNECTING ||
      currentState === WebSocket.OPEN
    ) {
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
      setStatus("disconnected");
      setMessages((prev) => [...prev, "Disconnected from server"]);
    };

    // クリーンアップ
    return () => {
      // 接続中または接続済みの場合のみclose()を呼ぶ
      if (
        websocket.readyState === WebSocket.CONNECTING ||
        websocket.readyState === WebSocket.OPEN
      ) {
        websocket.close();
      }
    };
  }, [url]);

  return {
    status,
    messages,
    sendMessage,
    getWebSocket,
    connect,
    disconnect,
  };
}
