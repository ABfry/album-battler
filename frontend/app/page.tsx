// app/title/page.tsx
"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useRoom } from "@/src/hooks/useRoom";
import { useWebSocket } from "@/src/lib/websocket/contexts/WebSocketContext";
import { NeedLoginButton } from "@/src/components/need-login-button/container/NeedLoginButton";

export default function TitlePage() {
  const router = useRouter();

  // WebSocket接続
  const { connect, disconnect } = useWebSocket();

  const {
    createRoom,
    joinRoom,
    leaveRoom,
    startGame,
    getBattleID,
    loading,
    error,
  } = useRoom();

  const handleCreateRoom = async () => {
    // 4桁のランダムな部屋番号を発行（例: 1111〜9999）
    const roomId = Math.floor(1000 + Math.random() * 9000).toString();
    const result = await createRoom("550e8400-e29b-41d4-a716-446655440001");
    router.push(`/room/${result?.room_id}`);
  };

  // WebSocketを接続
  useEffect(() => {
    connect();
  }, [connect, disconnect]);

  return (
    <main className="flex min-h-screen items-center justify-center">
      <div className="flex h-[640px] w-[360px] flex-col items-center justify-center gap-6 rounded-2xl border border-gray-400 bg-white shadow-lg">
        <h1 className="mb-6 text-xl font-semibold">アルバムバトラー</h1>

        <NeedLoginButton
          loggedInOnClick={handleCreateRoom}
          content="部屋をつくる"
        />

        <NeedLoginButton
          loggedInOnClick={() => router.push("/search")}
          content="部屋をさがす"
        />
      </div>
    </main>
  );
}
