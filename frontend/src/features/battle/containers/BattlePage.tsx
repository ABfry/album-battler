"use client";

import { useState, useRef, useEffect } from "react";
import { useBattleInfo } from "@/src/hooks/useBattleInfo";
import { useRoomInfo } from "@/src/hooks/useRoomInfo";
import { useWebSocketEvents } from "@/src/lib/websocket/hooks/useWebSocketEvents";
import type {
  ImageSendPayload,
  PlayerJoinRoomPayload,
} from "@/src/lib/websocket/types";
import { Battle } from "../components/Battle";

type BattlePageProps = {
  battleID: string;
};

/**
 * バトル画面のコンテナコンポーネント (Container)
 * ロジック・状態管理を担当
 */
export function BattlePage({ battleID }: BattlePageProps) {
  const [selectedImage, setSelectedImage] = useState<string | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [timeLeft, setTimeLeft] = useState(30);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const { subscribe } = useWebSocketEvents();

  // バトル情報の取得
  const {
    battle,
    images,
    loading: battleInfoLoading,
    error: battleInfoError,
    refetch: refetchBattle,
    refetchImages,
  } = useBattleInfo(battleID);

  // 部屋情報の取得
  const { room, refetch } = useRoomInfo(battle?.roomId || null);

  const players = room?.users;

  // WebSocketイベントの購読
  useEffect(() => {
    const unsubscribeJoin = subscribe(
      "player_join_room",
      (payload: PlayerJoinRoomPayload) => {
        console.log("Player joined room:", payload.room_id);
        refetch();
      }
    );
    const unsubscribeImageSend = subscribe(
      "image_send",
      (payload: ImageSendPayload) => {
        console.log("Image sent:", payload);
        // 画像送信イベントを受信したら画像一覧を再取得
        refetchImages();
      }
    );

    return () => {
      unsubscribeJoin();
      unsubscribeImageSend();
    };
  }, [subscribe, refetchImages, refetch]);

  // 画像選択ハンドラ
  const handleImageSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (event) => {
      const result = event.target?.result as string;
      setSelectedImage(result);
    };
    reader.readAsDataURL(file);
  };

  // アルバムを開くボタンのクリックハンドラ
  const handleOpenAlbum = () => {
    fileInputRef.current?.click();
  };

  // 画像選択を取り消すハンドラ
  const handleCancel = () => {
    setSelectedImage(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  // ドラッグアンドドロップハンドラ
  const handleDragEnter = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);

    const file = e.dataTransfer.files?.[0];
    if (!file || !file.type.startsWith("image/")) return;

    const reader = new FileReader();
    reader.onload = (event) => {
      const result = event.target?.result as string;
      setSelectedImage(result);
    };
    reader.readAsDataURL(file);
  };

  // クリップボードからの貼り付けハンドラ
  useEffect(() => {
    const handlePaste = (e: ClipboardEvent) => {
      const items = e.clipboardData?.items;
      if (!items) return;

      for (let i = 0; i < items.length; i++) {
        const item = items[i];
        if (item.type.startsWith("image/")) {
          const file = item.getAsFile();
          if (!file) continue;

          const reader = new FileReader();
          reader.onload = (event) => {
            const result = event.target?.result as string;
            setSelectedImage(result);
          };
          reader.readAsDataURL(file);
          break;
        }
      }
    };

    window.addEventListener("paste", handlePaste);
    return () => window.removeEventListener("paste", handlePaste);
  }, []);

  // カウントダウンタイマー
  useEffect(() => {
    if (timeLeft <= 0) return;

    const timer = setInterval(() => {
      setTimeLeft((prev) => Math.max(0, prev - 1));
    }, 1000);

    return () => clearInterval(timer);
  }, [timeLeft]);

  return (
    <Battle
      theme={battle?.theme || null}
      isLoading={battleInfoLoading}
      error={battleInfoError}
      players={players}
      images={images}
      selectedImage={selectedImage}
      isDragging={isDragging}
      timeLeft={timeLeft}
      fileInputRef={fileInputRef}
      onImageSelect={handleImageSelect}
      onOpenAlbum={handleOpenAlbum}
      onCancel={handleCancel}
      onDragEnter={handleDragEnter}
      onDragLeave={handleDragLeave}
      onDragOver={handleDragOver}
      onDrop={handleDrop}
    />
  );
}
