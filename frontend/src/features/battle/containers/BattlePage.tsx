"use client";

import { useState, useRef, useEffect } from "react";
import { useBattle } from "@/src/hooks/useBattle";
import { useBattleInfo } from "@/src/hooks/useBattleInfo";
import { useRoomInfo } from "@/src/hooks/useRoomInfo";
import { useWebSocketEvents } from "@/src/lib/websocket/hooks/useWebSocketEvents";
import type {
  ImageSendPayload,
  PlayerJoinRoomPayload,
  PlayerLeaveRoomPayload,
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
  const [sendImageSuccess, setSendImageSuccess] = useState(false);
  const [isImageSent, setIsImageSent] = useState(false); // 画像送信済みフラグ
  const fileInputRef = useRef<HTMLInputElement>(null);
  const { subscribe } = useWebSocketEvents();

  // バトル作成・画像送信の操作
  const {
    createBattle,
    sendImage,
    loading: battleLoading,
    error: battleError,
  } = useBattle();

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
        refetchBattle();
      }
    );
    const unsubscribeLeave = subscribe(
      "player_leave_room",
      (payload: PlayerLeaveRoomPayload) => {
        console.log("Player left room:", payload.room_id);
        refetch();
        refetchBattle();
      }
    );
    const unsubscribeImageSend = subscribe(
      "image_send",
      (payload: ImageSendPayload) => {
        console.log("Image sent:", payload);
        // 画像送信イベントを受信したら画像一覧を再取得
        refetchBattle();
        refetchImages();
      }
    );

    return () => {
      unsubscribeJoin();
      unsubscribeLeave();
      unsubscribeImageSend();
    };
  }, [subscribe, refetchBattle, refetchImages, refetch]);

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

  // 画像を確定して送信するハンドラ
  const handleConfirmImage = async () => {
    console.log("handleConfirmImage called");
    console.log("selectedImage:", selectedImage?.substring(0, 50));
    console.log("battleID:", battleID);

    if (!selectedImage) {
      console.log("No selected image, returning");
      return;
    }

    // Data URLからBase64部分のみを抽出
    // 例: "data:image/png;base64,iVBORw0KG..." -> "iVBORw0KG..."
    const base64 = selectedImage.split(",")[1];
    if (!base64) {
      console.error("Failed to extract base64 from data URL");
      return;
    }

    // TODO: 実際のユーザーIDを取得する仕組みが必要
    // 現在は仮のユーザーIDを使用
    const userId = "550e8400-e29b-41d4-a716-446655440001";

    console.log("Sending image with userId:", userId);
    console.log("Base64 length:", base64.length);
    const success = await sendImage(battleID, userId, base64);
    console.log("Send image result:", success);

    if (success) {
      console.log("Image sent successfully");
      // 送信成功時は画像をそのまま表示し、操作を無効化する
      setIsImageSent(true);
      setSendImageSuccess(true);
    } else {
      console.error("Failed to send image");
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
      isImageSent={isImageSent}
      fileInputRef={fileInputRef}
      onImageSelect={handleImageSelect}
      onOpenAlbum={handleOpenAlbum}
      onCancel={handleCancel}
      onConfirmImage={handleConfirmImage}
      onDragEnter={handleDragEnter}
      onDragLeave={handleDragLeave}
      onDragOver={handleDragOver}
      onDrop={handleDrop}
    />
  );
}
