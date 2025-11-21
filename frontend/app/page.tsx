// app/title/page.tsx
"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useRoom } from "@/src/hooks/useRoom";
import { useWebSocket } from "@/src/lib/websocket/contexts/WebSocketContext";
import { NeedLoginButton } from "@/src/components/need-login-button/container/NeedLoginButton";
import { getUserIdClient } from "@/src/lib/auth/getUserIdClient";

export default function TitlePage() {
  const router = useRouter();

  // WebSocket接続
  const { connect, disconnect } = useWebSocket();

  const { createRoom } = useRoom();

  const handleCreateRoom = async () => {
    const userId = getUserIdClient() || "";
    if (userId === "") {
      console.error("ユーザーIDがありません");
      return;
    }
    const result = await createRoom(getUserIdClient() || "");
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
