// app/search/page.tsx
"use client";

import { useState, type ChangeEvent } from "react";
import { useRouter } from "next/navigation";
import { useRoom } from "@/src/hooks/useRoom";
import { getUserIdClient } from "@/src/lib/auth/getUserIdClient";
import Link from "next/link";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/src/components/ui/alert-dialog";

export default function SearchPage() {
  const [roomId, setRoomId] = useState("");
  const [showErrorDialog, setShowErrorDialog] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");
  const router = useRouter();
  const { joinRoom } = useRoom();

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    const numbersOnly = value.replace(/[^\d]/g, "");
    const fourDigits = numbersOnly.match(/^\d{0,4}/)?.[0] ?? "";
    setRoomId(fourDigits);
  };

  const handleJoin = async () => {
    if (roomId.trim() === "") {
      setErrorMessage("部屋番号を入力してください");
      setShowErrorDialog(true);
      return;
    }
    if (roomId.length !== 4) {
      setErrorMessage("部屋番号は数字4桁で入力してください");
      setShowErrorDialog(true);
      return;
    }
    const userId = getUserIdClient() || "";
    const joinRoomId = await joinRoom(userId, Number(roomId));
    if (!joinRoomId) {
      setErrorMessage("指定された部屋が見つかりませんでした");
      setShowErrorDialog(true);
      return;
    }
    router.push(`/room/${joinRoomId}`);
  };

  return (
    <main className="flex min-h-screen items-center justify-center">
      <div className="flex w-full max-w-xs flex-col items-center gap-6 px-6 py-10">
        <h1 className="text-2xl font-black text-[#4a3b2a]">ルームに参加</h1>

        <input
          type="text"
          inputMode="numeric"
          value={roomId}
          onChange={handleChange}
          className="w-full rounded-xl border-2 border-[#c6c0b5] bg-white px-4 py-2 text-center text-xl font-semibold tracking-widest shadow-[0_3px_0_#b3ac9f] placeholder:text-[#b8b2a7] focus:border-[#a39c8e] focus:outline-none"
          placeholder="部屋番号"
        />

        <button
          onClick={handleJoin}
          className="w-56 rounded-xl bg-[#6b5337] py-3 text-lg font-black text-white shadow-[0_6px_0_rgba(0,0,0,0.35)] transition hover:brightness-110 active:translate-y-1 active:shadow-[0_2px_0_rgba(0,0,0,0.35)]"
        >
          参加
        </button>

        <Link //タイトルに戻るボタン
          href="/"
          className="absolute top-4 left-4 flex h-10 w-10 items-center justify-center rounded-md bg-white text-xl shadow"
        >
          ◀
        </Link>
      </div>

      <AlertDialog
        open={showErrorDialog}
        onOpenChange={(open) => setShowErrorDialog(open)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>入室エラー</AlertDialogTitle>
            <AlertDialogDescription>{errorMessage}</AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogAction onClick={() => setShowErrorDialog(false)}>
              OK
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </main>
  );
}
