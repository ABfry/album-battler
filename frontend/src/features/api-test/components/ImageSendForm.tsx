import { useState, useRef } from "react";
import { fileToBase64 } from "@/src/lib/utils";

type Props = {
  battleId: string | null;
  onSendImage: (userId: string, imageBase64: string) => void;
  loading: boolean;
  error: string | null;
  success: boolean;
};

/**
 * 画像送信フォーム (Presentational)
 */
export function ImageSendForm({
  battleId,
  onSendImage,
  loading,
  error,
  success,
}: Props) {
  const [imageBase64, setImageBase64] = useState("");
  const [fileName, setFileName] = useState("");
  const [userId, setUserId] = useState("550e8400-e29b-41d4-a716-446655440001");
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) {
      console.log("No file selected");
      return;
    }

    console.log("File selected:", file.name, file.type, file.size);

    // ファイル名を保存
    setFileName(file.name);

    try {
      // fileToBase64 utilを使用してBase64に変換
      const base64 = await fileToBase64(file);
      console.log("Base64 converted, length:", base64.length);
      setImageBase64(base64);
    } catch (error) {
      console.error("Failed to convert file to base64:", error);
      setFileName("");
    }
  };

  const handleSubmit = () => {
    if (imageBase64.trim() && userId.trim()) {
      onSendImage(userId, imageBase64);
    }
  };

  const handleClear = () => {
    setImageBase64("");
    setFileName("");
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  return (
    <div className="rounded-lg border p-4">
      <h2 className="mb-4 text-xl font-bold">9. Send Image</h2>

      <div className="mb-3 space-y-2 rounded bg-gray-50 p-3">
        <div>
          <p className="text-sm text-gray-600">Battle ID:</p>
          <p className="font-mono text-sm">
            {battleId || "No battle selected"}
          </p>
        </div>
        <div>
          <p className="text-sm text-gray-600">User ID:</p>
          <input
            type="text"
            value={userId}
            onChange={(e) => setUserId(e.target.value)}
            placeholder="User ID を入力"
            className="w-full rounded border bg-white px-2 py-1 font-mono text-sm"
          />
        </div>
        <div>
          <p className="mb-1 text-sm text-gray-600">画像を選択:</p>
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            onChange={handleFileSelect}
            className="hidden"
          />
          <button
            type="button"
            onClick={() => fileInputRef.current?.click()}
            className="w-full rounded border-2 border-dashed border-gray-300 bg-white px-4 py-3 text-sm hover:border-teal-500 hover:bg-teal-50"
          >
            {fileName ? (
              <span className="text-green-600">✓ {fileName}</span>
            ) : (
              <span className="text-gray-600">📁 ファイルを選択</span>
            )}
          </button>
        </div>
        {imageBase64 && (
          <div>
            <div className="mb-1 flex items-center justify-between">
              <p className="text-sm text-gray-600">
                Image Base64 (プレビュー):
              </p>
              <button
                onClick={handleClear}
                className="text-xs text-red-600 hover:underline"
              >
                クリア
              </button>
            </div>
            <textarea
              value={imageBase64.substring(0, 200) + "..."}
              readOnly
              className="w-full rounded border bg-gray-100 px-2 py-1 font-mono text-xs"
              rows={3}
            />
            <p className="mt-1 text-xs text-gray-500">
              長さ: {imageBase64.length} 文字
            </p>
          </div>
        )}
      </div>

      <button
        onClick={handleSubmit}
        disabled={!battleId || !imageBase64.trim() || loading}
        className="w-full rounded bg-teal-500 px-4 py-2 text-white hover:bg-teal-600 disabled:bg-gray-300"
      >
        {loading ? "Sending..." : "Send Image"}
      </button>

      {error && (
        <div className="mt-3 rounded bg-red-100 p-3 text-red-700">
          Error: {error}
        </div>
      )}

      {success && (
        <div className="mt-3 rounded bg-green-100 p-3 text-green-700">
          ✅ Image sent successfully!
        </div>
      )}
    </div>
  );
}
