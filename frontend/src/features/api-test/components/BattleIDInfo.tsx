type Props = {
  roomId: string | null;
  battleID: string | null;
  onGetBattleID: () => void;
  loading: boolean;
  error: string | null;
};

/**
 * バトルID取得・表示 (Presentational)
 */
export function BattleIDInfo({
  roomId,
  battleID,
  onGetBattleID,
  loading,
  error,
}: Props) {
  return (
    <div className="rounded-lg border p-4 shadow-sm">
      <h2 className="mb-4 text-xl font-bold">6. Get Battle ID</h2>

      <div className="mb-3 space-y-2 rounded bg-gray-50 p-3">
        <p className="text-sm text-gray-600">
          start_gameイベントを受け取って自動で走ります．
        </p>
        <div>
          <p className="text-sm text-gray-600">Room ID:</p>
          <p className="font-mono text-sm">{roomId || "No room selected"}</p>
        </div>
      </div>

      <button
        onClick={onGetBattleID}
        disabled={!roomId || loading}
        className="w-full rounded bg-purple-500 px-4 py-2 text-white hover:bg-purple-600 disabled:bg-gray-300"
      >
        {loading ? "Getting..." : "Get Battle ID"}
      </button>

      {error && (
        <div className="mt-3 rounded bg-red-100 p-3 text-red-700">
          Error: {error}
        </div>
      )}

      {battleID && (
        <div className="mt-3 rounded bg-green-50 p-3">
          <p className="text-sm font-semibold text-green-700">Battle ID:</p>
          <p className="font-mono text-xs break-all text-green-600">
            {battleID}
          </p>
        </div>
      )}
    </div>
  );
}
