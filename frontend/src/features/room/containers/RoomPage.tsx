"use client";

import { useWebSocket } from "@/app/providers/WebSocketProvider";
import { Room } from "../components/Room";

/**
 * RoomPageコンテナ
 * 実際のロジック
 */
export function RoomPage() {
  const { status, messages} = useWebSocket();

  return <Room status={status} messages={messages} connectionUrl="ws://localhost:8080/ws" />;
  // TODO: connectionUrlを環境変数から取得する、Providerも要修正?
}
