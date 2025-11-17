type Props = {
  roomId: string | null;
  userId: string;
  onStartGame: () => void;
  loading: boolean;
  error: string | null;
  success: boolean;
};

/**
 * ゲーム開始ボタン (Presentational)
 */
export function GameStartButton({
  roomId,
  userId,
  onStartGame,
  loading,
  error,
  success,
}: Props) {
  return (
    <div className="rounded-lg border bg-white p-4">
      <h2 className="mb-4 text-xl font-bold">5. Start Game</h2>

      <div className="mb-3 space-y-2 rounded bg-gray-50 p-3">
        <div>
          <p className="text-sm text-gray-600">Room ID:</p>
          <p className="font-mono text-sm">{roomId || "No room selected"}</p>
        </div>
        <div>
          <p className="text-sm text-gray-600">User ID (Host):</p>
          <input
            type="text"
            value={userId}
            readOnly
            className="w-full rounded border bg-white px-2 py-1 font-mono text-sm"
          />
        </div>
      </div>

      <button
        onClick={onStartGame}
        disabled={!roomId || loading}
        className="w-full rounded bg-purple-500 px-4 py-2 text-white hover:bg-purple-600 disabled:bg-gray-300"
      >
        {loading ? "Starting..." : "Start Game"}
      </button>

      {error && (
        <div className="mt-3 rounded bg-red-100 p-3 text-red-700">
          Error: {error}
        </div>
      )}

      {success && (
        <div className="mt-3 rounded bg-green-100 p-3 text-green-700">
          ✅ Game started!
        </div>
      )}
    </div>
  );
}
