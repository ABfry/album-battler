import { useEffect, useRef, useState } from "react";

export type ConnectionStatus =
  | "connecting"
  | "connected"
  | "disconnected"
  | "error";

export type WebSocketConnection = {
  status: ConnectionStatus;
  messages: string[];
  sendMessage: (data: string) => void;
  getWebSocket: () => WebSocket | null;
};

/**
 * WebSocket接続を管理するカスタムフック
 * @param url - WebSocketサーバのURL
 * @returns WebSocket接続の状態と操作関数
 */
export function useWebSocketConnection(url: string): WebSocketConnection {
  const wsRef = useRef<WebSocket | null>(null);
  const [status, setStatus] = useState<ConnectionStatus>("connecting");
  const [messages, setMessages] = useState<string[]>([]);

  // WebSocketインスタンスを取得する関数
  const getWebSocket = () => wsRef.current;

  // メッセージ送信関数
  const sendMessage = (data: string) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(data);
    } else {
      console.warn("WebSocket is not open. Cannot send message.");
    }
  };

  useEffect(() => {
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
      if (websocket.readyState === WebSocket.OPEN) {
        websocket.close();
      }
    };
  }, [url]);

  return {
    status,
    messages,
    sendMessage,
    getWebSocket,
  };
}
