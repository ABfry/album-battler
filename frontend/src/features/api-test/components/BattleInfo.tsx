import Link from "next/link";
import type { Battle } from "@/src/lib/api/types";

type Props = {
  battle: Battle | null;
  loading: boolean;
  error: string | null;
  onRefetch: () => void;
};

/**
 * バトル情報表示 (Presentational)
 */
export function BattleInfo({ battle, loading, error, onRefetch }: Props) {
  return (
    <div className="rounded-lg border bg-white p-4">
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-xl font-bold">8. Battle Info</h2>
        <div className="flex gap-2">
          {battle && (
            <Link
              href={`/battle/${battle.id}`}
              className="rounded bg-green-500 px-3 py-1 text-sm text-white hover:bg-green-600"
            >
              Open Battle
            </Link>
          )}
          <button
            onClick={onRefetch}
            disabled={loading}
            className="rounded bg-blue-500 px-3 py-1 text-sm text-white hover:bg-blue-600 disabled:bg-gray-300"
          >
            {loading ? "Loading..." : "Refetch"}
          </button>
        </div>
      </div>

      <div className="mb-3 space-y-2 rounded bg-gray-50 p-3">
        <p className="text-sm text-gray-600">
          💡 battleIDセット時に自動でfetchされます
        </p>
      </div>

      {error && (
        <div className="rounded bg-red-100 p-3 text-red-700">
          Error: {error}
        </div>
      )}

      {battle ? (
        <div className="space-y-2 text-sm">
          <div>
            <p className="font-semibold text-gray-600">Battle ID:</p>
            <p className="font-mono text-xs break-all">{battle.id}</p>
          </div>
          <div>
            <p className="font-semibold text-gray-600">Room ID:</p>
            <p className="font-mono text-xs break-all">{battle.roomId}</p>
          </div>
          <div>
            <p className="font-semibold text-gray-600">Theme:</p>
            <p>{battle.theme}</p>
          </div>
          <div>
            <p className="font-semibold text-gray-600">Started At:</p>
            <p className="text-xs">
              {new Date(battle.startedAt).toLocaleString()}
            </p>
          </div>
          <div>
            <p className="font-semibold text-gray-600">
              Participants ({battle.userIds.length}):
            </p>
            <ul className="mt-1 space-y-1">
              {battle.userIds.map((userId, idx) => (
                <li key={userId} className="font-mono text-xs text-gray-600">
                  {idx + 1}. {userId}
                </li>
              ))}
            </ul>
          </div>
        </div>
      ) : (
        <p className="text-sm text-gray-500">No battle data</p>
      )}
    </div>
  );
}
