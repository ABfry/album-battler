// app/title/page.tsx
"use client";

import { useRouter } from "next/navigation";
import { useRoom } from "@/src/hooks/useRoom";

export default function TitlePage() {
  const router = useRouter();

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

  return (
    <main className="flex min-h-screen items-center justify-center">
      <div className="flex h-[640px] w-[360px] flex-col items-center justify-center gap-6 rounded-2xl border border-gray-400 bg-white shadow-lg">
        <h1 className="mb-6 text-xl font-semibold">アルバムバトラー</h1>

        <button
          onClick={handleCreateRoom}
          className="w-40 rounded-md border border-gray-500 bg-white py-2 shadow-sm transition hover:-translate-y-0.5 hover:shadow"
        >
          部屋をつくる
        </button>

        <button
          onClick={() => router.push("/search")}
          className="w-40 rounded-md border border-gray-500 bg-white py-2 shadow-sm transition hover:-translate-y-0.5 hover:shadow"
        >
          部屋をさがす
        </button>

        <button
          onClick={() => router.push("/watch")}
          className="w-40 rounded-md border border-gray-500 bg-white py-2 shadow-sm transition hover:-translate-y-0.5 hover:shadow"
        >
          観戦？
        </button>
      </div>
    </main>
  );
}
