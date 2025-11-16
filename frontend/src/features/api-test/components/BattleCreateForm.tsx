type Props = {
  roomId: string | null;
  onCreateBattle: () => void;
  loading: boolean;
  error: string | null;
  result: { BattleID: string } | null;
};

/**
 * バトル作成フォーム (Presentational)
 */
export function BattleCreateForm({
  roomId,
  onCreateBattle,
  loading,
  error,
  result,
}: Props) {
  return (
    <div className="rounded-lg border p-4">
      <h2 className="mb-4 text-xl font-bold">7. Create Battle</h2>

      <div className="mb-3 space-y-2 rounded bg-gray-50 p-3">
        <div>
          <p className="text-sm text-gray-600">Room ID:</p>
          <p className="font-mono text-sm">{roomId || "No room selected"}</p>
        </div>
      </div>

      <button
        onClick={onCreateBattle}
        disabled={!roomId || loading}
        className="w-full rounded bg-indigo-500 px-4 py-2 text-white hover:bg-indigo-600 disabled:bg-gray-300"
      >
        {loading ? "Creating..." : "Create Battle"}
      </button>

      {error && (
        <div className="mt-3 rounded bg-red-100 p-3 text-red-700">
          Error: {error}
        </div>
      )}

      {result && (
        <div className="mt-3 rounded bg-green-100 p-3">
          <p className="text-sm font-semibold text-green-700">
            Battle Created!
          </p>
          <p className="font-mono text-xs break-all text-green-600">
            {result.BattleID}
          </p>
        </div>
      )}
    </div>
  );
}
