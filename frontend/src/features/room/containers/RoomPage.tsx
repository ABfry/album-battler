"use client";

import { useEffect } from "react";
import { useWebSocket } from "@/src/lib/websocket/contexts/WebSocketContext";
import { Room } from "../components/Room";

/**
 * RoomPageコンテナ
 * 実際のロジック
 */
export function RoomPage() {
  const { status, messages, connect } = useWebSocket();

  const connectionUrl =
    process.env.NEXT_PUBLIC_WEBSOCKET_URL || "ws://localhost:8080/ws";

  // マウント時にWebSocket接続を開始
  useEffect(() => {
    connect();

    // 接続は Provider がアンマウントされるまで維持（ページ遷移では切断しない）
  }, [connect]);

  return (
    <Room status={status} messages={messages} connectionUrl={connectionUrl} />
  );
}
