"use client";

import {
  createContext,
  useContext,
  useEffect,
  useRef,
  useState,
  ReactNode,
} from "react";

type WebSocketContextType = {
  getWebSocket: () => WebSocket | null;
  status: "connecting" | "connected" | "disconnected" | "error";
  messages: string[];
};

const WebSocketContext = createContext<WebSocketContextType>({
  getWebSocket: () => null,
  status: "connecting",
  messages: [],
});

export function useWebSocket() {
  return useContext(WebSocketContext);
}

type WebSocketProviderProps = {
  children: ReactNode;
};

export function WebSocketProvider({ children }: WebSocketProviderProps) {
  const wsRef = useRef<WebSocket | null>(null);
  const [status, setStatus] = useState<
    "connecting" | "connected" | "disconnected" | "error"
  >("connecting");
  const [messages, setMessages] = useState<string[]>([]);

  // WebSocketインスタンスを取得する関数
  const getWebSocket = () => wsRef.current;

  useEffect(() => {
    // WebSocket接続を作成
    const websocket = new WebSocket("ws://localhost:8080/ws");
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
  }, []);

  return (
    <WebSocketContext.Provider value={{ getWebSocket, status, messages }}>
      {children}
    </WebSocketContext.Provider>
  );
}
