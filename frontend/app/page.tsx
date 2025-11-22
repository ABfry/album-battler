// app/title/page.tsx
"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useRoom } from "@/src/hooks/useRoom";
import { useWebSocket } from "@/src/lib/websocket/contexts/WebSocketContext";
import { NeedLoginButton } from "@/src/components/need-login-button/container/NeedLoginButton";
import { getUserIdClient } from "@/src/lib/auth/getUserIdClient";

import TitleBackSlider from "@/src/features/title/components/BackSlider";
import { TitleLogo } from "@/src/components/ui/title-logo";

export default function TitlePage() {
  const router = useRouter();

  // WebSocket接続（常に呼び出す、URLが空の場合は接続しない）
  const { connect } = useWebSocket();

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

  // userIdがある場合のみWebSocketを接続
  useEffect(() => {
    const userId = getUserIdClient();
    if (userId) {
      connect();
    }
  }, [connect]);

  return (
    <main className="relative flex min-h-screen flex-col items-center justify-center">
      <TitleLogo className="mx-auto w-4/5 md:w-96" />
      <div className="flex flex-col items-center justify-center">
        <NeedLoginButton
          loggedInOnClick={handleCreateRoom}
          onLoginSuccess={() => window.location.reload()}
          content="部屋をつくる"
          className="px-12 py-5 text-xl"
        />

        <NeedLoginButton
          loggedInOnClick={() => router.push("/search")}
          onLoginSuccess={() => window.location.reload()}
          content="部屋をさがす"
          className="px-12 py-5 text-xl"
        />
      </div>
      <TitleBackSlider />
    </main>
  );
}
