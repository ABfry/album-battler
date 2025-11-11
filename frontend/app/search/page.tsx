// app/search/page.tsx
"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

export default function SearchPage() {
  const [roomId, setRoomId] = useState("");
  const router = useRouter();

  const handleJoin = () => {
    const trimmed = roomId.trim();
    if (!trimmed) {
      alert("部屋番号を入力してください");
      return;
    }
    router.push(`/room/${trimmed}`);
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
            value={roomId}
            onChange={(e) => setRoomId(e.target.value)}
            className="w-full border-b border-gray-500 pb-1 text-center focus:outline-none"
            placeholder="0000"
            maxLength={4}
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
