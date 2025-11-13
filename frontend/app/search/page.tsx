// app/search/page.tsx
"use client";

import { useState, type ChangeEvent } from "react";
import { useRouter } from "next/navigation";

export default function SearchPage() {
  const [roomId, setRoomId] = useState("");
  const router = useRouter();

  // ルームID入力のハンドラ（数字のみ・最大4桁）
  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;

    // 1. 数字以外をすべて削除
    const numbersOnly = value.replace(/[^\d]/g, "");

    // 2. 先頭から最大4桁だけを取り出す
    const fourDigits = numbersOnly.match(/^\d{0,4}/)?.[0] ?? "";

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
    <main className="flex min-h-screen items-center justify-center bg-gray-100">
      <div className="flex h-[640px] w-[360px] flex-col items-center rounded-2xl border border-gray-400 bg-white px-6 pt-12 shadow-lg">
        <h1 className="mb-24 text-xl font-semibold">アルバムバトラー</h1>

        <div className="flex w-full max-w-xs flex-col items-center gap-4 rounded-xl border-2 border-gray-400 px-6 py-10">
          <p className="text-center text-sm leading-relaxed">
            部屋番号を
            <br />
            入力してください
          </p>

          <input
            type="text"
            inputMode="numeric" // スマホで数字キーボードを出す
            value={roomId}
            onChange={handleChange}
            className="w-full border-b border-gray-500 pb-1 text-center focus:outline-none"
            placeholder="1234"
          />

          <button
            onClick={handleJoin}
            className="mt-4 rounded-md border border-gray-500 px-8 py-2 text-sm shadow-sm transition hover:-translate-y-0.5 hover:shadow"
          >
            参加
          </button>
        </div>

        <button
          onClick={() => router.push("/title")}
          className="mt-10 text-xs text-gray-500 hover:underline"
        >
          ← タイトルに戻る
        </button>
      </div>
    </main>
  );
}
