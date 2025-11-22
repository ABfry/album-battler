// app/room/[roomId]/page.tsx
"use client";

import { useParams } from "next/navigation";
import { useEffect, useState } from "react";
import { useRoomInfo } from "@/src/hooks/useRoomInfo";
import type {
  PlayerJoinRoomPayload,
  StartGamePayload,
  RoomSettingsUpdatedPayload,
} from "@/src/lib/websocket/types";
import { useWebSocketEvents } from "@/src/lib/websocket/hooks/useWebSocketEvents";
import Link from "next/link";
import { useRoom } from "@/src/hooks/useRoom";
import { useRouter } from "next/navigation";
import { Button } from "@/src/components/ui/button";
import { getUserIdClient } from "@/src/lib/auth/getUserIdClient";

export default function RoomPage() {
  const { roomID } = useParams() as { roomID: string };

  const { startGame, getBattleID, updateRoomSettings } = useRoom();

  const router = useRouter();

  // 部屋情報の取得（作成後に自動取得）
  const { room, refetch } = useRoomInfo(roomID);
  const { subscribe } = useWebSocketEvents();
  const [isUpdating, setIsUpdating] = useState(false);

  const userId = getUserIdClient() || "";
  const isHost = room?.hostUserId === userId;
  useEffect(() => {
    const unsubscribeJoin = subscribe(
      "player_join_room",
      (payload: PlayerJoinRoomPayload) => {
        console.log("Player joined room:", payload.room_id);
        refetch();
      }
    );

    const unsubscribeStart = subscribe(
      "start_game",
      async (payload: StartGamePayload) => {
        console.log("Game started:", payload.room_id);
        refetch();
        const id = await getBattleID(roomID);
        router.push(`/battle/${id}`);
      }
    );

    const unsubscribeSettings = subscribe(
      "room_settings_updated",
      (payload: RoomSettingsUpdatedPayload) => {
        console.log(
          "Room settings updated:",
          payload.battle_time_limit_seconds
        );
        refetch();
      }
    );

    return () => {
      unsubscribeJoin();
      unsubscribeStart();
      unsubscribeSettings();
    };
  }, [subscribe, refetch, getBattleID, roomID, router]);

  const handleBattle = async () => {
    const userId = getUserIdClient() || "";
    if (userId === "") {
      console.error("ユーザーIDがありません");
      return;
    }
    await startGame(roomID, userId);

    const battleId = await getBattleID(roomID);

    router.push(`/battle/${battleId}`);
  };

  const handleTimeLimitChange = async (
    e: React.ChangeEvent<HTMLSelectElement>
  ) => {
    const newTimeLimit = parseInt(e.target.value, 10);
    if (!userId || !isHost) return;

    setIsUpdating(true);
    try {
      const success = await updateRoomSettings(roomID, userId, {
        battle_time_limit_seconds: newTimeLimit,
      });
      if (success) {
        await refetch();
      } else {
        console.error("Failed to update time limit");
      }
    } catch (error) {
      console.error("Failed to update time limit:", error);
    } finally {
      setIsUpdating(false);
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

        {/* バトル時間設定・表示 */}
        <div className="mb-4 w-full max-w-xs rounded-lg border border-[#3551b8] bg-white px-4 py-3">
          <label className="mb-2 block text-sm font-bold text-gray-700">
            バトル時間
          </label>
          {isHost ? (
            <select
              value={room?.battleTimeLimitSeconds ?? 60}
              onChange={handleTimeLimitChange}
              disabled={isUpdating}
              className="w-full rounded border border-gray-300 bg-white px-3 py-2 text-base font-medium disabled:opacity-50"
            >
              <option value={30}>30秒</option>
              <option value={60}>1分</option>
              <option value={120}>2分</option>
              <option value={180}>3分</option>
              <option value={240}>4分</option>
              <option value={300}>5分</option>
            </select>
          ) : (
            <div className="w-full rounded border border-gray-300 bg-gray-50 px-3 py-2 text-base font-medium text-gray-700">
              {room?.battleTimeLimitSeconds === 30 && "30秒"}
              {room?.battleTimeLimitSeconds === 60 && "1分"}
              {room?.battleTimeLimitSeconds === 120 && "2分"}
              {room?.battleTimeLimitSeconds === 180 && "3分"}
              {room?.battleTimeLimitSeconds === 240 && "4分"}
              {room?.battleTimeLimitSeconds === 300 && "5分"}
            </div>
          )}
        </div>

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
          disabled={(room?.users?.length ?? 0) < 2}
          onClick={handleBattle}
          className="w-full"
        >
          バトル！
        </Button>
      </div>
    </main>
  );
}
