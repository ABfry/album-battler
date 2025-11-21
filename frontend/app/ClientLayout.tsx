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
  const wsUrl = userId
    ? `${process.env.NEXT_PUBLIC_WEBSOCKET_URL || "ws://localhost:8080/ws"}?user_id=${userId}`
    : null;

  // ユーザーIDがない場合はWebSocket接続なしで子要素をレンダリング
  if (!wsUrl) {
    return <>{children}</>;
  }

  return <WebSocketProvider url={wsUrl}>{children}</WebSocketProvider>;
}
