// app/room/[roomId]/page.tsx
"use client";

import { useParams, useRouter } from "next/navigation";

export default function RoomDetailPage() {
  // URLのパラメータ（例: /room/1111 → roomID = "1111"）
  const { roomID } = useParams() as { roomID: string };
  const router = useRouter();

  // 仮のプレイヤーデータ
  const players = ["井上", "岩崎"];
  const watcherCount = 5;

  const handleReady = () => {
    alert("準備完了の処理を書く（APIなど）");
  };

  const handleExit = () => {
    router.push("/title");
  };

  return (
    <main className="flex min-h-screen items-center justify-center bg-gray-100">
      <div className="flex h-[640px] w-[360px] flex-col rounded-2xl border border-gray-400 bg-white px-6 py-10 shadow-lg">
        {/* 部屋番号表示 */}
        <h2 className="mb-8 text-center text-xl font-semibold">
          部屋番号: {roomID}
        </h2>

        {/* プレイヤー一覧 */}
        <div className="mb-6">
          <p className="mb-1 font-semibold">プレイヤー</p>
          <ul className="list-inside list-disc text-sm">
            {players.map((name) => (
              <li key={name}>{name}</li>
            ))}
          </ul>
        </div>

        {/* 観戦者数 */}
        <div className="mb-10">
          <p className="mb-1 font-semibold">観戦</p>
          <ul className="list-inside list-disc text-sm">
            <li>{watcherCount}名</li>
          </ul>
        </div>

        {/* ボタン群 */}
        <div className="mt-auto flex flex-col items-center gap-4">
          <button
            onClick={handleReady}
            className="w-40 rounded-md border border-gray-500 bg-white py-2 shadow-sm transition hover:-translate-y-0.5 hover:shadow"
          >
            準備完了
          </button>
          <button
            onClick={handleExit}
            className="w-40 rounded-md border border-gray-500 bg-white py-2 shadow-sm transition hover:-translate-y-0.5 hover:shadow"
          >
            退出
          </button>
        </div>
      </div>
    </main>
  );
}
