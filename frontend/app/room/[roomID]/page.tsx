// app/room/[roomId]/page.tsx
"use client";

import { useParams, useRouter } from "next/navigation";
import { useState } from "react";

type Player = {
  id: number;
  name: string;
  joined: boolean;
};

export default function RoomPage() {
  const { roomId } = useParams() as { roomId: string };
  const router = useRouter();

  // 仮のプレイヤー情報（joined=false が「待機中..」枠）
  const [players] = useState<Player[]>([
    { id: 1, name: "岩崎", joined: true },
    { id: 2, name: "井上", joined: true },
    { id: 3, name: "シバタ", joined: true },
    { id: 4, name: "なかむら", joined: true },
    { id: 5, name: "待機中・・", joined: false },
  ]);

  const handleBattle = () => {
    alert("バトル開始の処理を書く");
  };

  return (
    <main className="flex min-h-screen items-center justify-center bg-[#d6c2a4]">
      {/* 戻るボタン */}
      <button
        onClick={() => router.push("/title")}
        className="absolute top-4 left-4 flex h-10 w-10 items-center justify-center rounded-md bg-white text-xl shadow"
      >
        ◀
      </button>

      {/* ルーム全体コンテナ（縦長スマホ想定） */}
      <div className="flex h-[640px] w-[360px] flex-col items-center">
        {/* タイトル */}
        <h1 className="mt-12 mb-4 text-3xl font-black tracking-widest text-[#b57c39]">
          ルーム
        </h1>

        {/* 部屋番号（必要なら表示） */}
        <p className="mb-4 text-xs text-gray-700">部屋番号: {roomId}</p>

        {/* プレイヤー一覧カード */}
        <div className="w-full max-w-xs rounded-xl border border-[#3551b8] bg-white px-6 py-6 shadow-[0_8px_0_rgba(0,0,0,0.15)]">
          <ul className="space-y-3">
            {players.map((p) => (
              <li key={p.id} className="flex items-center gap-3">
                {/* アイコンの丸 */}
                <div className="flex h-9 w-9 items-center justify-center rounded-full border border-black">
                  {/* 中の顔アイコンはシンプルに線だけ */}
                  <div className="h-5 w-5 rounded-full border border-gray-400" />
                </div>

                {/* 名前 */}
                <span
                  className={
                    p.joined
                      ? "text-lg font-black"
                      : "text-lg font-black text-gray-300"
                  }
                >
                  {p.name}
                </span>
              </li>
            ))}
          </ul>
        </div>

        {/* バトルボタン */}
        <button
          onClick={handleBattle}
          className="mt-8 w-56 rounded-xl bg-[#6b5337] py-3 text-lg font-black text-white shadow-[0_6px_0_rgba(0,0,0,0.35)] active:translate-y-1 active:shadow-[0_2px_0_rgba(0,0,0,0.35)]"
        >
          バトル！！
        </button>
      </div>
    </main>
  );
}
