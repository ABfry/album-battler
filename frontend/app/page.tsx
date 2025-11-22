// app/title/page.tsx
"use client";

import { useEffect, useCallback } from "react";
import { useRouter } from "next/navigation";
import { useRoom } from "@/src/hooks/useRoom";
import { useWebSocket } from "@/src/lib/websocket/contexts/WebSocketContext";
import { NeedLoginButton } from "@/src/components/need-login-button/container/NeedLoginButton";
import { getUserIdClient } from "@/src/lib/auth/getUserIdClient";

export default function TitlePage() {
  const router = useRouter();

  // WebSocket接続（常に呼び出す、URLが空の場合は接続しない）
  const { connect, setUrl } = useWebSocket();

  const { createRoom } = useRoom();

  // ログイン成功時: WebSocket URL を設定して接続
  const handleLoginSuccess = useCallback(() => {
    const userId = getUserIdClient();
    if (userId) {
      const wsUrl = `${process.env.NEXT_PUBLIC_WEBSOCKET_URL || "ws://localhost:8080/ws"}?user_id=${userId}`;
      setUrl(wsUrl);
      connect();
    }
  }, [setUrl, connect]);

  const handleCreateRoom = async () => {
    const userId = getUserIdClient() || "";
    if (userId === "") {
      console.error("ユーザーIDがありません");
      return;
    }
    const result = await createRoom(getUserIdClient() || "");
    router.push(`/room/${result?.room_id}`);
  };

  // userIdがある場合のみWebSocketを接続
  useEffect(() => {
    const userId = getUserIdClient();
    if (userId) {
      connect();
    }
  }, [connect]);

  return (
    <main className="flex min-h-screen items-center justify-center">
      <div className="flex h-[640px] w-[360px] flex-col items-center justify-center gap-6 rounded-2xl border border-gray-400 bg-white shadow-lg">
        <h1 className="mb-6 text-xl font-semibold">アルバムバトラー</h1>

        <NeedLoginButton
          loggedInOnClick={handleCreateRoom}
          onLoginSuccess={handleLoginSuccess}
          content="部屋をつくる"
        />

        <NeedLoginButton
          loggedInOnClick={() => router.push("/search")}
          onLoginSuccess={handleLoginSuccess}
          content="部屋をさがす"
        />
      </div>
    </main>
  );
}
