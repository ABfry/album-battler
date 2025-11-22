// app/room/[roomId]/page.tsx
"use client";

import { useParams } from "next/navigation";
import { useEffect, useState } from "react";
import { useRoomInfo } from "@/src/hooks/useRoomInfo";
import type {
  PlayerJoinRoomPayload,
  StartButtonPressedPayload,
  StartGamePayload,
} from "@/src/lib/websocket/types";
import { useWebSocketEvents } from "@/src/lib/websocket/hooks/useWebSocketEvents";
import Link from "next/link";
import { useRoom } from "@/src/hooks/useRoom";
import { useRouter } from "next/navigation";
import { Button } from "@/src/components/ui/button";
import { getUserIdClient } from "@/src/lib/auth/getUserIdClient";

export default function RoomPage() {
  const { roomID } = useParams() as { roomID: string };

  const { startGame, getBattleID } = useRoom();

  const router = useRouter();

  // ゲーム開始待機中のローディング状態
  const [isLoading, setIsLoading] = useState(false);

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

    const unsubscribePressed = subscribe(
      "start_button_pressed",
      async (payload: StartButtonPressedPayload) => {
        console.log("Start button Pressed:", payload.room_id);
        setIsLoading(true);
        refetch();
      }
    );

    const unsubscribeStart = subscribe(
      "start_game",
      async (payload: StartGamePayload) => {
        console.log("Game started:", payload.room_id);
        setIsLoading(false);
        refetch();
        const id = await getBattleID(roomID);
        router.push(`/battle/${id}`);
      }
    );
    return () => {
      unsubscribeJoin();
      unsubscribePressed();
      unsubscribeStart();
    };
  }, [subscribe, refetch, getBattleID, roomID, router]);

  const handleBattle = async () => {
    const userId = getUserIdClient() || "";
    if (userId === "") {
      console.error("ユーザーIDがありません");
      return;
    }
    setIsLoading(true);
    try {
      await startGame(roomID, userId);
    } catch {
      console.error("ゲーム開始エラー");
      setIsLoading(false);
    }
  };

  return (
    <main className="flex min-h-screen items-center justify-center bg-[#d6c2a4]">
      {/* 戻るボタン */}
      <Link
        href="/"
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
        <p className="mb-4 text-xs text-gray-700">
          部屋番号: {room?.roomNumber}
        </p>

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
        <Button
          variant="primary"
          size="lg"
          disabled={(room?.users?.length ?? 0) < 2 || isLoading}
          onClick={handleBattle}
          className="mt-8 w-full"
        >
          {isLoading ? (
            <div className="flex items-center justify-center gap-3">
              <div className="h-5 w-5 animate-spin rounded-full border-2 border-white/60 border-t-white" />
              <span className="text-lg font-black">待機中...</span>
            </div>
          ) : (
            <span className="text-lg font-black">バトル！</span>
          )}
        </Button>
      </div>
    </main>
  );
}
