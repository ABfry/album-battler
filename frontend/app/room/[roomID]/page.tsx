// app/title/page.tsx
"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useRoom } from "@/src/hooks/useRoom";
import { useWebSocket } from "@/src/lib/websocket/contexts/WebSocketContext";
import { NeedLoginButton } from "@/src/components/need-login-button/container/NeedLoginButton";
import { getUserIdClient } from "@/src/lib/auth/getUserIdClient";

export default function TitlePage() {
  const router = useRouter();

  // WebSocket接続（常に呼び出す、URLが空の場合は接続しない）
  const { connect } = useWebSocket();

  const { createRoom } = useRoom();

  // ★ 追加：loading状態
  const [loading, setLoading] = useState(false);

  const handleCreateRoom = async () => {
    const userId = getUserIdClient() || "";
    if (userId === "") {
      console.error("ユーザーIDがありません");
      return;
    }

    // ★ 追加：ローディング開始
    setLoading(true);

    const result = await createRoom(userId);

    // ★ 追加：ローディング終了
    setLoading(false);

    if (result?.room_id) {
      router.push(`/room/${result.room_id}`);
    }
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

        {/* ★ 変更：loading中は文言切り替え＆押せない */}
        <Button
          loggedInOnClick={handleCreateRoom}
          content={loading ? "作成中..." : "部屋をつくる"}
          disabled={loading}
        />

        <Button
          loggedInOnClick={() => router.push("/search")}
          content="部屋をさがす"
        />
      </div>
    </main>
  );
}
