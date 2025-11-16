"use client";

import { useState } from "react";
import Image from "next/image";

// モックデータ
const mockPlayers = [
  {
    id: 1,
    name: "井上",
    avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=test1",
  },
  {
    id: 2,
    name: "岩﨑",
    avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=test2",
  },
  {
    id: 3,
    name: "シバタ",
    avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=test3",
  },
  {
    id: 4,
    name: "なかむら",
    avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=test1",
  },
  {
    id: 5,
    name: "木村",
    avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=test2",
  },
];

const mockTheme = "謎の儀式現場";

export default function BattlePage() {
  const [players] = useState(mockPlayers);

  return (
    <div className="flex min-h-screen flex-col items-center justify-center p-4">
      <div className="flex h-screen w-full max-w-4xl flex-col items-center justify-between py-8">
        {/* お題表示（最上部） */}
        <div className="w-full text-center">
          <h1 className="text-4xl font-black text-slate-800 md:text-5xl">
            {mockTheme}
          </h1>
        </div>

        {/* バトル画像（中央） */}
        <div className="flex w-full flex-1 items-center justify-center">
          <div className="relative w-2/3 max-w-md">
            {/* 額縁の外枠 */}
            <div className="rounded-lg bg-gradient-to-br from-amber-800 via-amber-700 to-amber-900 p-4 shadow-2xl">
              {/* 額縁の内側（金色の装飾） */}
              <div className="rounded-md border-4 border-amber-600 bg-gradient-to-br from-amber-200 to-amber-300 p-3 shadow-inner">
                {/* 白いマット（正方形の固定サイズ） */}
                <div className="relative aspect-square w-full overflow-hidden rounded-sm border-2 border-amber-100 bg-white p-6 shadow-md">
                  <Image
                    src="/takoyaki.jpg"
                    alt="バトル画像"
                    fill
                    sizes="(max-width: 768px) 100vw, 400px"
                    className="object-contain"
                    priority
                  />
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* プレイヤー情報（最下部） */}
        <div className="w-full">
          {/* アイコン */}
          <div className="mb-4 flex items-center justify-center gap-4 md:gap-6">
            {players.map((player) => (
              <div key={player.id} className="flex flex-col items-center gap-2">
                {/* アイコン画像 */}
                <div className="h-16 w-16 overflow-hidden rounded-full border-4 border-slate-300 bg-slate-200 shadow-lg md:h-20 md:w-20">
                  <Image
                    src={player.avatar}
                    alt={player.name}
                    width={80}
                    height={80}
                    className="h-full w-full object-cover"
                  />
                </div>
              </div>
            ))}
          </div>

          {/* ユーザー名 */}
          <div className="flex items-center justify-center gap-4 md:gap-6">
            {players.map((player) => (
              <div key={player.id} className="w-16 text-center md:w-20">
                <p className="truncate text-sm font-bold text-slate-700 md:text-base">
                  {player.name}
                </p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
