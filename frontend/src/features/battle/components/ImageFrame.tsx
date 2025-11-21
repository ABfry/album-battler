import Image from "next/image";
import { Button } from "@/src/components/ui/button";

type ImageFrameProps = {
  selectedImage: string | null;
  displayedImage?: string | null; // 拍手フェーズで表示する画像
  isDragging: boolean;
  isImageSent: boolean;
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
  const imageToShow = displayedImage || selectedImage;
  return (
    <div className="flex w-full flex-1 flex-col items-center justify-center gap-4">
      <div className="relative w-2/3 max-w-md">
        {/* 隠しファイル入力 */}
        <input
          ref={fileInputRef}
          type="file"
          accept="image/*"
          onChange={onImageSelect}
          className="hidden"
          disabled={isImageSent}
        />

        {/* 額縁の外枠 */}
        <div className="rounded-lg bg-linear-to-br from-amber-800 via-amber-700 to-amber-900 p-4 shadow-2xl">
          {/* 額縁の内側（金色の装飾） */}
          <div className="rounded-md border-4 border-amber-600 bg-linear-to-br from-amber-200 to-amber-300 p-3 shadow-inner">
            {/* 白いマット（正方形の固定サイズ） */}
            <div
              className={`relative aspect-square w-full overflow-hidden rounded-sm border-2 bg-white p-6 shadow-md transition-colors ${
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
                  isImageSent || displayedImage ? undefined : onOpenAlbum
                }
                disabled={isImageSent || !!displayedImage}
                className={`flex h-full w-full items-center justify-center transition-colors ${
                  isImageSent || displayedImage
                    ? "cursor-not-allowed"
                    : "hover:bg-gray-50"
                }`}
              >
                {imageToShow ? (
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
                    <div className="mb-2 text-6xl">
                      {isDragging ? "📥" : "📂"}
                    </div>
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
      </div>

      {/* ボタン */}
      {selectedImage && !isImageSent && (
        <div className="flex gap-3">
          <Button
            variant="secondary"
            size="lg"
            onClick={() => {
              console.log("Cancel button clicked");
              onCancel();
            }}
          >
            取り消す
          </Button>
          <Button
            variant="primary"
            size="lg"
            onClick={() => {
              console.log("Confirm button clicked");
              onConfirmImage();
            }}
          >
            これで決定
          </Button>
        </div>
      )}
    </div>
  );
}
