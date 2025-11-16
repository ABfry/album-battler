import Image from "next/image";
import { Button } from "@/src/components/ui/button";

type ImageFrameProps = {
  selectedImage: string | null;
  isDragging: boolean;
  fileInputRef: React.RefObject<HTMLInputElement | null>;
  onImageSelect: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onOpenAlbum: () => void;
  onCancel: () => void;
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
  isDragging,
  fileInputRef,
  onImageSelect,
  onOpenAlbum,
  onCancel,
  onDragEnter,
  onDragLeave,
  onDragOver,
  onDrop,
}: ImageFrameProps) {
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
        />

        {/* 額縁の外枠 */}
        <div className="rounded-lg bg-linear-to-br from-amber-800 via-amber-700 to-amber-900 p-4 shadow-2xl">
          {/* 額縁の内側（金色の装飾） */}
          <div className="rounded-md border-4 border-amber-600 bg-linear-to-br from-amber-200 to-amber-300 p-3 shadow-inner">
            {/* 白いマット（正方形の固定サイズ） */}
            <div
              className={`relative aspect-square w-full overflow-hidden rounded-sm border-2 bg-white p-6 shadow-md transition-colors ${
                isDragging ? "border-blue-400 bg-blue-50" : "border-amber-100"
              }`}
              onDragEnter={onDragEnter}
              onDragLeave={onDragLeave}
              onDragOver={onDragOver}
              onDrop={onDrop}
            >
              {selectedImage ? (
                <Image
                  src={selectedImage}
                  alt="選択した画像"
                  fill
                  sizes="(max-width: 768px) 100vw, 400px"
                  className="object-contain"
                  unoptimized
                />
              ) : (
                <button
                  onClick={onOpenAlbum}
                  className="flex h-full w-full items-center justify-center transition-colors hover:bg-gray-50"
                >
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
                </button>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* 選びなおすボタン */}
      {selectedImage && (
        <div className="flex gap-3">
          <Button variant="primary" size="lg" onClick={onOpenAlbum}>
            選びなおす
          </Button>
          <Button variant="secondary" size="lg" onClick={onCancel}>
            取り消す
          </Button>
        </div>
      )}
    </div>
  );
}
