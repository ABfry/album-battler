import { useState } from "react";

type Props = {
  battleId: string | null;
  onSendClap: (
    userId: string,
    targetUserId: string,
    battleId: string,
    count: number
  ) => void;
  loading: boolean;
  success: boolean;
};

/**
 * 拍手送信フォーム (Presentational)
 */
export function ClapSendForm({
  battleId,
  onSendClap,
  loading,
  success,
}: Props) {
  const [userId, setUserId] = useState("550e8400-e29b-41d4-a716-446655440002");
  const [targetUserId, setTargetUserId] = useState(
    "550e8400-e29b-41d4-a716-446655440001"
  );
  const [count, setCount] = useState(5);

  const handleSubmit = () => {
    if (battleId && userId.trim() && targetUserId.trim()) {
      onSendClap(userId, targetUserId, battleId, count);
    }
  };

  return (
    <div className="rounded-lg border bg-white p-4">
      <h2 className="mb-4 text-xl font-bold">11. Send Clap</h2>

      <div className="mb-3 space-y-2 rounded bg-gray-50 p-3">
        <div>
          <p className="text-sm text-gray-600">Battle ID:</p>
          <p className="font-mono text-sm">
            {battleId || "No battle selected"}
          </p>
        </div>
        <div>
          <p className="text-sm text-gray-600">User ID (拍手する人):</p>
          <input
            type="text"
            value={userId}
            onChange={(e) => setUserId(e.target.value)}
            placeholder="User ID を入力"
            className="w-full rounded border bg-white px-2 py-1 font-mono text-sm"
          />
        </div>
        <div>
          <p className="text-sm text-gray-600">
            Target User ID (拍手される人):
          </p>
          <input
            type="text"
            value={targetUserId}
            onChange={(e) => setTargetUserId(e.target.value)}
            placeholder="Target User ID を入力"
            className="w-full rounded border bg-white px-2 py-1 font-mono text-sm"
          />
        </div>
        <div>
          <p className="mb-1 text-sm text-gray-600">Clap Count (1以上):</p>
          <input
            type="number"
            min="1"
            value={count}
            onChange={(e) => setCount(Number(e.target.value))}
            className="w-full rounded border bg-white px-2 py-1 font-mono text-sm"
          />
        </div>
      </div>

      <button
        onClick={handleSubmit}
        disabled={
          !battleId || !userId.trim() || !targetUserId.trim() || loading
        }
        className="w-full rounded bg-purple-500 px-4 py-2 text-white hover:bg-purple-600 disabled:bg-gray-300"
      >
        {loading ? "Sending..." : "Send Clap"}
      </button>

      {success && (
        <div className="mt-3 rounded bg-green-100 p-3 text-green-700">
          ✅ Clap sent successfully!
        </div>
      )}
    </div>
  );
}
