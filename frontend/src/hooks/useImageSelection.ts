import { useState, useRef, useCallback, useEffect } from "react";
import { useBattle } from "./useBattle";
import { compressImage } from "@/src/lib/utils";

type UseImageSelectionOptions = {
  battleId: string;
  userId: string;
  canSelect: boolean; // フェーズから受け取る
  onSendSuccess?: () => void;
  onSendError?: () => void;
};

// バックエンドが受け付ける画像形式
const ALLOWED_IMAGE_TYPES = [
  "image/jpeg",
  "image/png",
  "image/webp",
  "image/heic",
] as const;

/**
 * 画像形式が許可されているかチェック
 */
const isAllowedImageType = (type: string): boolean => {
  return ALLOWED_IMAGE_TYPES.includes(
    type as (typeof ALLOWED_IMAGE_TYPES)[number]
  );
};

/**
 * 画像選択ロジック管理フック
 * ファイル入力、ドラッグ&ドロップ、クリップボード対応
 * フェーズに応じた有効/無効制御
 */
export function useImageSelection(options: UseImageSelectionOptions) {
  const [selectedImage, setSelectedImage] = useState<string | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [isImageSent, setIsImageSent] = useState(false);
  const [isCompressing, setIsCompressing] = useState(false);
  const [isSending, setIsSending] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const { sendImage } = useBattle();

  // 共通の画像処理ロジック
  const processImageFile = useCallback(async (file: File) => {
    // 画像形式チェック
    if (!isAllowedImageType(file.type)) {
      alert(
        `サポートされていない画像形式です。\nJPEG、PNG、WebP、HEICのいずれかを選択してください。`
      );
      return;
    }

    try {
      setIsCompressing(true);
      const compressedDataUrl = await compressImage(file);
      setSelectedImage(compressedDataUrl);
    } catch (error) {
      console.error("[useImageSelection] 圧縮エラー:", error);
      alert("画像の圧縮に失敗しました。別の画像を選択してください。");
    } finally {
      setIsCompressing(false);
    }
  }, []);

  // 画像選択ハンドラー
  const handleImageSelect = useCallback(
    async (e: React.ChangeEvent<HTMLInputElement>) => {
      if (!options.canSelect || isImageSent) return;

      const file = e.target.files?.[0];
      if (!file) return;

      await processImageFile(file);
    },
    [options.canSelect, isImageSent, processImageFile]
  );

  // アルバムを開く
  const handleOpenAlbum = useCallback(() => {
    if (!options.canSelect || isImageSent) return;
    fileInputRef.current?.click();
  }, [options.canSelect, isImageSent]);

  // 画像選択を取り消す
  const handleCancel = useCallback(() => {
    if (isImageSent) return;
    setSelectedImage(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  }, [isImageSent]);

  // 画像を確定して送信
  const handleConfirmImage = useCallback(async () => {
    console.log("[useImageSelection] handleConfirmImage called");
    console.log(
      "[useImageSelection] selectedImage:",
      selectedImage?.substring(0, 50)
    );
    console.log("[useImageSelection] battleId:", options.battleId);
    console.log("[useImageSelection] canSelect:", options.canSelect);

    if (!options.canSelect || !selectedImage || isImageSent || isSending) {
      console.log(
        "[useImageSelection] Cannot confirm: canSelect=%s, hasImage=%s, alreadySent=%s, isSending=%s",
        options.canSelect,
        !!selectedImage,
        isImageSent,
        isSending
      );
      return;
    }

    // Data URLからBase64部分のみを抽出
    // 例: "data:image/png;base64,iVBORw0KG..." -> "iVBORw0KG..."
    const base64 = selectedImage.split(",")[1];
    if (!base64) {
      console.error(
        "[useImageSelection] Failed to extract base64 from data URL"
      );
      return;
    }

    setIsSending(true);

    try {
      console.log(
        "[useImageSelection] Sending image with userId:",
        options.userId
      );
      console.log("[useImageSelection] Base64 length:", base64.length);
      const success = await sendImage(options.battleId, options.userId, base64);
      console.log("[useImageSelection] Send image result:", success);

      if (success) {
        console.log("[useImageSelection] Image sent successfully");
        // 送信成功時は画像をそのまま表示し、操作を無効化する
        setIsImageSent(true);
        options.onSendSuccess?.();
      } else {
        console.error("[useImageSelection] Failed to send image");
        // 送信失敗時はエラーダイアログを表示
        options.onSendError?.();
      }
    } finally {
      setIsSending(false);
    }
  }, [selectedImage, isImageSent, isSending, options, sendImage]);

  // ドラッグ&ドロップハンドラー
  const handleDragEnter = useCallback(
    (e: React.DragEvent) => {
      if (!options.canSelect || isImageSent) return;
      e.preventDefault();
      e.stopPropagation();
      setIsDragging(true);
    },
    [options.canSelect, isImageSent]
  );

  const handleDragLeave = useCallback(
    (e: React.DragEvent) => {
      if (!options.canSelect || isImageSent) return;
      e.preventDefault();
      e.stopPropagation();
      setIsDragging(false);
    },
    [options.canSelect, isImageSent]
  );

  const handleDragOver = useCallback(
    (e: React.DragEvent) => {
      if (!options.canSelect || isImageSent) return;
      e.preventDefault();
      e.stopPropagation();
    },
    [options.canSelect, isImageSent]
  );

  const handleDrop = useCallback(
    async (e: React.DragEvent) => {
      if (!options.canSelect || isImageSent) return;
      e.preventDefault();
      e.stopPropagation();
      setIsDragging(false);

      const file = e.dataTransfer.files?.[0];
      if (!file || !file.type.startsWith("image/")) return;

      await processImageFile(file);
    },
    [options.canSelect, isImageSent, processImageFile]
  );

  // クリップボードからの貼り付けハンドラー
  useEffect(() => {
    if (!options.canSelect || isImageSent) return;

    const handlePaste = async (e: ClipboardEvent) => {
      const items = e.clipboardData?.items;
      if (!items) return;

      for (let i = 0; i < items.length; i++) {
        const item = items[i];
        if (item.type.startsWith("image/")) {
          const file = item.getAsFile();
          if (!file) continue;

          await processImageFile(file);
          break;
        }
      }
    };

    window.addEventListener("paste", handlePaste);
    return () => window.removeEventListener("paste", handlePaste);
  }, [options.canSelect, isImageSent, processImageFile]);

  return {
    selectedImage,
    isDragging,
    isImageSent,
    isSending,
    isCompressing,
    fileInputRef,
    handleImageSelect,
    handleOpenAlbum,
    handleCancel,
    handleConfirmImage,
    handleDragEnter,
    handleDragLeave,
    handleDragOver,
    handleDrop,
  };
}
