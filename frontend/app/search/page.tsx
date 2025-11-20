// app/search/page.tsx
"use client";

import { useState, type ChangeEvent } from "react";
import { useRouter } from "next/navigation";

export default function SearchPage() {
  const [roomId, setRoomId] = useState("");
  const router = useRouter();

  // ルームID入力（数字のみ・最大4桁）
  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;

    // 全角数字 → 半角数字に変換
    const normalized = value.replace(/[０-９]/g, (ch) =>
      String.fromCharCode(ch.charCodeAt(0) - 0xfee0)
    );

    // 半角数字以外を削除
    const onlyNumbers = normalized.replace(/[^0-9]/g, "");

    // 先頭から最大4桁
    const fourDigits = onlyNumbers.match(/^[0-9]{0,4}/)?.[0] ?? "";

    setRoomId(fourDigits);
  };

  const handleJoin = () => {
    if (roomId.length !== 4) {
      alert("部屋番号は数字4桁で入力してください");
      return;
    }
    router.push(`/room/${roomId}`);
  };

  return (
    <main className="flex min-h-screen items-center justify-center bg-[#d6c2a4]">
      {/* 戻るボタン（RoomPage と同じUI） */}
      <button
        onClick={() => router.push("/title")}
        className="absolute top-4 left-4 flex h-10 w-10 items-center justify-center rounded-md bg-white text-xl shadow"
      >
        ◀
      </button>

      {/* 縦長スマホコンテナ（RoomPage と同じサイズ） */}
      <div className="flex h-[640px] w-[360px] flex-col items-center">
        {/* タイトル（雰囲気を揃える） */}
        <h1 className="mt-12 mb-4 text-3xl font-black tracking-widest text-[#b57c39]">
          ルーム検索
        </h1>

        {/* 説明テキスト */}
        <p className="mb-4 text-center text-xs text-gray-700">
          部屋番号を入力して
          <br />
          ルームに参加しよう
        </p>

        {/* 入力カード（枠・影を RoomPage のプレイヤーカードに揃え） */}
        <div className="w-full max-w-xs rounded-xl border border-[#3551b8] bg-white px-6 py-6 shadow-[0_8px_0_rgba(0,0,0,0.15)]">
          <p className="mb-4 text-center text-sm leading-relaxed">
            部屋番号（4桁）を
            <br />
            入力してください
          </p>

          <input
            type="text"
            inputMode="numeric"
            value={roomId}
            onChange={handleChange}
            className="w-full border-b border-gray-500 pb-1 text-center text-lg tracking-widest focus:outline-none"
            placeholder="1234"
          />

          {/* 参加ボタン：バトルボタンと同じUI */}
          <div className="mt-6 flex justify-center">
            <button
              onClick={handleJoin}
              className="w-56 rounded-xl bg-[#6b5337] py-3 text-lg font-black text-white shadow-[0_6px_0_rgba(0,0,0,0.35)] active:translate-y-1 active:shadow-[0_2px_0_rgba(0,0,0,0.35)]"
            >
              参加
            </button>
          </div>
        </div>
      </div>
    </main>
  );
}
