// app/room/[roomId]/page.tsx
"use client";

import { useParams } from "next/navigation";
import { useEffect } from "react";
import { useRoomInfo } from "@/src/hooks/useRoomInfo";
import type { PlayerJoinRoomPayload } from "@/src/lib/websocket/types";
import { useWebSocketEvents } from "@/src/lib/websocket/hooks/useWebSocketEvents";
import Link from "next/link";

export default function RoomPage() {
  const { roomID } = useParams() as { roomID: string };

  // 部屋情報の取得（作成後に自動取得）
  const { room, refetch } = useRoomInfo(roomID);

  const { subscribe } = useWebSocketEvents();
  useEffect(() => {
    const unsubscribeJoin = subscribe(
      "player_join_room",
      (payload: PlayerJoinRoomPayload) => {
        console.log("Player joined room:", payload.room_id);
        refetch();
      }
    );
    return () => {
      unsubscribeJoin();
    };
  }, [subscribe, refetch]);

  const handleBattle = () => {
    alert("バトル開始の処理を書く");
  };

  return (
    <main className="flex min-h-screen items-center justify-center bg-[#d6c2a4]">
      {/* 戻るボタン */}
      <Link
        href="/title"
        className="absolute top-4 left-4 flex h-10 w-10 items-center justify-center rounded-md bg-white text-xl shadow"
      >
        ◀
      </Link>

      {/* ルーム全体コンテナ（縦長スマホ想定） */}
      <div className="flex h-[640px] w-[360px] flex-col items-center">
        {/* タイトル */}
        <h1 className="mt-12 mb-4 text-3xl font-black tracking-widest text-[#b57c39]">
          ルーム
        </h1>

        {/* 部屋番号（必要なら表示） */}
        <p className="mb-4 text-xs text-gray-700">部屋番号: {roomID}</p>

        {/* プレイヤー一覧カード */}
        <div className="w-full max-w-xs rounded-xl border border-[#3551b8] bg-white px-6 py-6 shadow-[0_8px_0_rgba(0,0,0,0.15)]">
          <ul className="space-y-3">
            {room?.users.map((p) => (
              <li key={p.id} className="flex items-center gap-3">
                {/* アイコンの丸 */}
                <div className="flex h-9 w-9 items-center justify-center rounded-full border border-black">
                  {/* 中の顔アイコンはシンプルに線だけ */}
                  <div className="h-5 w-5 rounded-full border border-gray-400" />
                </div>

                {/* 名前 */}
                <span className={"text-lg font-black"}>{p.name}</span>
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
