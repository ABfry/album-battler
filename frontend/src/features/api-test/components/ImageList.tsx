import Image from "next/image";
import type { BattleImage } from "@/src/lib/api/types";

type Props = {
  images: BattleImage[];
  loading: boolean;
  error: string | null;
  onRefetch: () => void;
};

/**
 * バトル画像一覧表示 (Presentational)
 */
export function ImageList({ images, loading, error, onRefetch }: Props) {
  return (
    <div className="rounded-lg border p-4">
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-xl font-bold">10. Battle Images</h2>
        <button
          onClick={onRefetch}
          disabled={loading}
          className="rounded bg-blue-500 px-3 py-1 text-sm text-white hover:bg-blue-600 disabled:bg-gray-300"
        >
          {loading ? "Loading..." : "Refetch"}
        </button>
      </div>

      <div className="mb-3 space-y-2 rounded bg-gray-50 p-3">
        <p className="text-sm text-gray-600">
          💡 image_sendイベント受信時に自動でfetchされます
        </p>
      </div>

      {error && (
        <div className="rounded bg-red-100 p-3 text-red-700">
          Error: {error}
        </div>
      )}

      {images.length > 0 ? (
        <div className="space-y-4">
          {images.map((image, idx) => (
            <div key={idx} className="rounded border bg-gray-50 p-3">
              <div className="mb-2">
                <p className="text-xs text-gray-600">User ID:</p>
                <p className="font-mono text-xs break-all">{image.userId}</p>
              </div>
              <div className="mb-2">
                <p className="text-xs text-gray-600">Image URL:</p>
                <p className="font-mono text-xs break-all text-blue-600">
                  {image.imageUrl}
                </p>
              </div>
              <div>
                <p className="mb-1 text-xs text-gray-600">画像:</p>
                <div className="relative h-64 w-full">
                  <Image
                    src={image.imageUrl}
                    alt={`User ${image.userId}'s image`}
                    fill
                    className="rounded border object-contain"
                    unoptimized
                  />
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <p className="text-sm text-gray-500">No images yet</p>
      )}
    </div>
  );
}
