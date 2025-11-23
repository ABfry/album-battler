import Image from "next/image";
import { motion, AnimatePresence } from "framer-motion";
import { Button } from "@/src/components/ui/button";

type ImageFrameProps = {
  selectedImage: string | null;
  displayedImage?: string | null; // 拍手フェーズで表示する画像
  isDragging: boolean;
  isImageSent: boolean;
  isSending?: boolean;
  isCompressing?: boolean;
  canSelect: boolean; // 画像選択可能か（タイムアップ判定用）
  fileInputRef: React.RefObject<HTMLInputElement | null>;
  onImageSelect: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onOpenAlbum: () => void;
  onCancel: () => void;
  onConfirmImage: () => void;
  onDragEnter: (e: React.DragEvent) => void;
  onDragLeave: (e: React.DragEvent) => void;
  onDragOver: (e: React.DragEvent) => void;
  onDrop: (e: React.DragEvent) => void;
};

/**
 * 画像選択用の額縁コンポーネント (Presentational)
 */
export function ImageFrame({
  selectedImage,
  displayedImage,
  isDragging,
  isImageSent,
  isSending = false,
  isCompressing = false,
  canSelect,
  fileInputRef,
  onImageSelect,
  onOpenAlbum,
  onCancel,
  onConfirmImage,
  onDragEnter,
  onDragLeave,
  onDragOver,
  onDrop,
}: ImageFrameProps) {
  // 表示する画像を決定: displayedImage > selectedImage
  // displayedImage が null の場合は画像未提出として扱う
  const isDisplayingOtherPlayer = displayedImage !== undefined;
  const isNoImageSubmitted = displayedImage === null;
  const imageToShow = displayedImage || selectedImage;

  // 額縁コンテンツを共通化
  const frameContent = (
    <div className="h-full w-full rounded-lg bg-linear-to-br from-amber-800 via-amber-700 to-amber-900 p-4 shadow-2xl">
      {/* 額縁の内側（金色の装飾） */}
      <div className="h-full w-full rounded-md border-4 border-amber-600 bg-linear-to-br from-amber-200 to-amber-300 p-3 shadow-inner">
        {/* 白いマット（正方形を維持） */}
        <div
          className={`relative h-full w-full overflow-hidden rounded-sm border-2 bg-white p-6 shadow-md transition-colors ${
            isDragging && !isImageSent
              ? "border-blue-400 bg-blue-50"
              : "border-amber-100"
          }`}
          onDragEnter={isImageSent ? undefined : onDragEnter}
          onDragLeave={isImageSent ? undefined : onDragLeave}
          onDragOver={isImageSent ? undefined : onDragOver}
          onDrop={isImageSent ? undefined : onDrop}
        >
          <button
            onClick={
              isImageSent || displayedImage || isCompressing
                ? undefined
                : onOpenAlbum
            }
            disabled={isImageSent || !!displayedImage || isCompressing}
            className={`flex h-full w-full items-center justify-center transition-colors ${
              isImageSent || displayedImage || isCompressing
                ? "cursor-not-allowed"
                : "hover:bg-gray-50"
            }`}
          >
            {isCompressing ? (
              <div className="text-center">
                <div className="mb-4 text-6xl">⏳</div>
                <p className="text-lg font-semibold text-gray-700">
                  画像を圧縮中...
                </p>
                <p className="mt-2 text-sm text-gray-500">
                  しばらくお待ちください
                </p>
              </div>
            ) : isNoImageSubmitted ? (
              <div className="text-center">
                <div className="mb-2 text-6xl">🚫</div>
                <p className="text-lg font-semibold text-gray-500">
                  画像未提出
                </p>
              </div>
            ) : imageToShow ? (
              <Image
                src={imageToShow}
                alt="選択した画像"
                fill
                sizes="(max-width: 768px) 100vw, 400px"
                className="object-contain"
                unoptimized
              />
            ) : (
              <div className="text-center">
                <div className="mb-2 text-6xl">{isDragging ? "📥" : "📂"}</div>
                <p className="text-lg font-semibold text-gray-700">
                  {isDragging ? "ここにドロップ" : "アルバムを開く"}
                </p>
                {!isDragging && (
                  <p className="mt-2 text-sm text-gray-500">
                    ドラッグ&ドロップ または Ctrl+V
                  </p>
                )}
              </div>
            )}
          </button>
        </div>
      </div>
    </div>
  );

  return (
    <div className="flex h-full w-full flex-col items-center justify-center gap-4">
      {/* 隠しファイル入力 */}
      <input
        ref={fileInputRef}
        type="file"
        accept="image/jpeg,image/png,image/webp,image/heic"
        onChange={onImageSelect}
        className="hidden"
        disabled={isImageSent}
      />

      {/* 額縁コンテナ: 利用可能なスペース内で最大の正方形を維持 */}
      <div
        className="relative min-h-0 w-full flex-1"
        style={{ containerType: "size" }}
      >
        <div className="absolute inset-0 flex items-center justify-center">
          {isDisplayingOtherPlayer ? (
            <AnimatePresence mode="wait">
              <motion.div
                key={displayedImage ?? "no-image"}
                initial={{ x: "100%", opacity: 0 }}
                animate={{ x: 0, opacity: 1 }}
                exit={{ x: "-100%", opacity: 0 }}
                transition={{
                  duration: 0.5,
                  ease: "easeInOut",
                }}
                style={{
                  aspectRatio: "1 / 1",
                  width: "min(100%, 100cqmin)",
                  height: "auto",
                  maxHeight: "100%",
                }}
              >
                {frameContent}
              </motion.div>
            </AnimatePresence>
          ) : (
            <div
              style={{
                aspectRatio: "1 / 1",
                width: "100%",
                height: "100%",
                maxWidth: "min(100cqw, 100cqh)",
                maxHeight: "min(100cqw, 100cqh)",
              }}
            >
              {frameContent}
            </div>
          )}
        </div>
      </div>

      {/* ボタン or 時間切れ表示 */}
      {selectedImage && !isImageSent && (
        <div className="flex w-full shrink-0 flex-wrap justify-center gap-1 px-2 sm:gap-3">
          {canSelect ? (
            <>
              <Button
                variant="secondary"
                size="lg"
                disabled={isSending}
                onClick={() => {
                  console.log("Cancel button clicked");
                  onCancel();
                }}
                className="mt-2 w-36 sm:mt-8 sm:w-56"
              >
                取り消す
              </Button>
              <Button
                variant="primary"
                size="lg"
                disabled={isSending}
                onClick={() => {
                  console.log("Confirm button clicked");
                  onConfirmImage();
                }}
                className="mt-2 w-36 sm:mt-8 sm:w-56"
              >
                {isSending ? "送信中..." : "これで決定"}
              </Button>
            </>
          ) : (
            <p className="text-lg font-bold text-red-500">時間切れ！</p>
          )}
        </div>
      )}
    </div>
  );
}
