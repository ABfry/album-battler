"use client";

import { useWebSocket } from "@/src/lib/websocket/contexts/WebSocketContext";
import { Room } from "../components/Room";

/**
 * RoomPageコンテナ
 * 実際のロジック
 */
export function RoomPage() {
  const { status, messages } = useWebSocket();

  const connectionUrl =
    process.env.NEXT_PUBLIC_WEBSOCKET_URL || "ws://localhost:8080/ws";

  return (
    <Room status={status} messages={messages} connectionUrl={connectionUrl} />
  );
}
