// app/room/[roomId]/page.tsx
"use client";

import { useParams } from "next/navigation";
import Image from "next/image";
import { useEffect, useState } from "react";
import { useRoomInfo } from "@/src/hooks/useRoomInfo";
import type {
  PlayerJoinRoomPayload,
  StartButtonPressedPayload,
  StartGamePayload,
  RoomSettingsUpdatedPayload,
} from "@/src/lib/websocket/types";
import { useWebSocketEvents } from "@/src/lib/websocket/hooks/useWebSocketEvents";
import Link from "next/link";
import { useRoom } from "@/src/hooks/useRoom";
import { useRouter } from "next/navigation";
import { Button } from "@/src/components/ui/button";
import { getUserIdClient } from "@/src/lib/auth/getUserIdClient";
import { Loading } from "@/src/components/ui/loading";

export default function RoomPage() {
  const { roomID } = useParams() as { roomID: string };

  const { startGame, getBattleID, updateRoomSettings } = useRoom();

  const router = useRouter();

  // ゲーム開始待機中のローディング状態
  const [isLoading, setIsLoading] = useState(false);

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
      unsubscribePressed();
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
    setIsLoading(true);
    try {
      await startGame(roomID, userId);
    } catch {
      console.error("ゲーム開始エラー");
      setIsLoading(false);
    }
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
    <main className="flex min-h-screen items-center justify-center">
      {/* ルーム全体コンテナ */}
      <div className="flex w-full max-w-sm flex-col items-center gap-6 px-6 py-10">
        {/* タイトル */}
        <h1 className="text-2xl font-black text-[#4a3b2a]">ルーム</h1>

        {/* 部屋番号 */}
        <div className="rounded-xl bg-white px-6 py-3 shadow-[0_3px_0_#b3ac9f]">
          <p className="text-center text-lg font-bold tracking-widest text-[#4a3b2a]">
            部屋番号: {room?.roomNumber}
          </p>
        </div>

        {/* バトル時間設定・表示 */}
        <div className="w-full rounded-xl bg-white px-6 py-4 shadow-[0_3px_0_#b3ac9f]">
          <label className="mb-3 block text-center text-sm font-bold text-[#4a3b2a]">
            バトル時間
          </label>
          {isHost ? (
            <select
              value={room?.battleTimeLimitSeconds ?? 60}
              onChange={handleTimeLimitChange}
              disabled={isUpdating}
              className="w-full rounded-xl border-2 border-[#c6c0b5] bg-white px-4 py-2 text-center text-base font-semibold focus:border-[#a39c8e] focus:outline-none disabled:opacity-50"
            >
              <option value={30}>30秒</option>
              <option value={60}>1分</option>
              <option value={120}>2分</option>
              <option value={180}>3分</option>
              <option value={240}>4分</option>
              <option value={300}>5分</option>
            </select>
          ) : (
            <div className="w-full rounded-xl border-2 border-[#c6c0b5] bg-gray-50 px-4 py-2 text-center text-base font-semibold text-[#4a3b2a]">
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
        <div className="w-full rounded-xl bg-white px-6 py-6 shadow-[0_3px_0_#b3ac9f]">
          <h2 className="mb-4 text-center text-sm font-bold text-[#4a3b2a]">
            参加プレイヤー
          </h2>
          <ul className="space-y-3">
            {room?.users.map((p) => (
              <li key={p.id} className="flex items-center gap-3">
                {/* アイコンの丸 */}
                <div className="flex h-9 w-9 items-center justify-center overflow-hidden rounded-full border-2 border-[#c6c0b5]">
                  <Image
                    width={36}
                    height={36}
                    src={p.icon_url}
                    alt={p.name}
                    className="h-full w-full rounded-full object-cover"
                  />
                </div>

                {/* 名前 */}
                <span className="text-lg font-bold text-[#4a3b2a]">
                  {p.name}
                </span>
              </li>
            ))}
          </ul>
        </div>

        {/* バトルボタン */}
        <Button
          disabled={!isHost || (room?.users?.length ?? 0) < 2 || isLoading}
          onClick={handleBattle}
          className="w-64 rounded-xl bg-[#6b5337] py-3 text-lg font-black text-white shadow-[0_6px_0_rgba(0,0,0,0.35)] transition hover:bg-[#6b5337] hover:brightness-110 active:translate-y-1 active:shadow-[0_2px_0_rgba(0,0,0,0.35)]"
        >
          {isLoading ? (
            <div className="flex items-center justify-center gap-3">
              <Loading />
              <span>待機中...</span>
            </div>
          ) : isHost ? (
            "バトル！"
          ) : (
            "ホストの開始を待っています"
          )}
        </Button>

        {/* 戻るボタン */}
        <Link
          href="/"
          className="absolute top-4 left-4 flex h-10 w-10 items-center justify-center rounded-md bg-white text-xl shadow"
        >
          ◀
        </Link>
      </div>
    </main>
  );
}
