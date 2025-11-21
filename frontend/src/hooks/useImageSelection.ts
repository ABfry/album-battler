import { useState, useRef, useCallback, useEffect } from "react";
import { useBattle } from "./useBattle";

type UseImageSelectionOptions = {
  battleId: string;
  userId: string;
  canSelect: boolean; // フェーズから受け取る
  onSendSuccess?: () => void;
  onSendError?: () => void;
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
  const fileInputRef = useRef<HTMLInputElement>(null);
  const { sendImage } = useBattle();

  // 画像選択ハンドラー
  const handleImageSelect = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      if (!options.canSelect || isImageSent) return;

      const file = e.target.files?.[0];
      if (!file) return;

      const reader = new FileReader();
      reader.onload = (event) => {
        const result = event.target?.result as string;
        setSelectedImage(result);
      };
      reader.readAsDataURL(file);
    },
    [options.canSelect, isImageSent]
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

    if (!selectedImage || isImageSent) {
      console.log(
        "[useImageSelection] No selected image or already sent, returning"
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
  }, [selectedImage, isImageSent, options, sendImage]);

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
    (e: React.DragEvent) => {
      if (!options.canSelect || isImageSent) return;
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
    },
    [options.canSelect, isImageSent]
  );

  // クリップボードからの貼り付けハンドラー
  useEffect(() => {
    if (!options.canSelect || isImageSent) return;

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
  }, [options.canSelect, isImageSent]);

  return {
    selectedImage,
    isDragging,
    isImageSent,
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
