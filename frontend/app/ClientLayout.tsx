"use client";

import { ReactNode } from "react";
import { WebSocketProvider } from "@/src/lib/websocket/providers/WebSocketProvider";

type ClientLayoutProps = {
  children: ReactNode;
  userId: string | null;
};

/**
 * WebSocket接続を管理するClient Component
 * ユーザーIDが存在する場合のみWebSocket接続を確立する
 */
export function ClientLayout({ children, userId }: ClientLayoutProps) {
  // ユーザーIDが存在する場合のみWebSocket URLを構築
  // userIdがない場合は空文字列（接続しない）
  const wsUrl = userId
    ? `${process.env.NEXT_PUBLIC_WEBSOCKET_URL || "ws://localhost:8080/ws"}?user_id=${userId}`
    : "";

  // 常にWebSocketProviderで囲む（urlが空文字列の場合は接続しない）
  return <WebSocketProvider url={wsUrl}>{children}</WebSocketProvider>;
}
