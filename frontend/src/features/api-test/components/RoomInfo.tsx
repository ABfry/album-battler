import Link from "next/link";
import Image from "next/image";
import type { Room } from "@/src/lib/api/types";

type Props = {
  room: Room | null;
  loading: boolean;
  error: string | null;
  onRefetch: () => void;
};

/**
 * 部屋情報表示 (Presentational)
 */
export function RoomInfo({ room, loading, error, onRefetch }: Props) {
  return (
    <div className="rounded-lg border p-4">
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-xl font-bold">3. Room Info</h2>
        <div className="flex gap-2">
          {room && (
            <Link
              href={`/room/${room.id}`}
              className="rounded bg-green-500 px-3 py-1 text-sm text-white hover:bg-green-600"
            >
              Open Room
            </Link>
          )}
          <button
            onClick={onRefetch}
            disabled={loading}
            className="rounded bg-gray-500 px-3 py-1 text-sm text-white hover:bg-gray-600 disabled:bg-gray-300"
          >
            {loading ? "Refreshing..." : "Refresh"}
          </button>
        </div>
      </div>

      {error && (
        <div className="rounded bg-red-100 p-3 text-red-700">
          Error: {error}
        </div>
      )}

      {loading && <p>Loading...</p>}

      {room && (
        <div className="space-y-2">
          <div className="rounded bg-gray-50 p-3">
            <p className="text-sm text-gray-600">Room Number</p>
            <p className="text-lg font-semibold">{room.roomNumber}</p>
          </div>

          <div className="rounded bg-gray-50 p-3">
            <p className="text-sm text-gray-600">Status</p>
            <p className="text-lg font-semibold">{room.status}</p>
          </div>

          <div className="rounded bg-gray-50 p-3">
            <p className="text-sm text-gray-600">Host</p>
            <p className="font-mono text-sm">{room.hostUserId || "No host"}</p>
          </div>

          <div className="rounded bg-gray-50 p-3">
            <p className="mb-2 text-sm text-gray-600">
              Users ({room.users.length})
            </p>
            <ul className="space-y-2">
              {room.users.map((user) => (
                <li
                  key={user.id}
                  className="flex items-center gap-2 rounded bg-white p-2"
                >
                  <Image
                    src={user.icon_url}
                    alt={user.name}
                    width={32}
                    height={32}
                    className="h-8 w-8 rounded-full"
                  />
                  <div>
                    <p className="font-medium">{user.name}</p>
                    <p className="text-xs text-gray-500">{user.id}</p>
                  </div>
                </li>
              ))}
            </ul>
          </div>

          <div className="rounded bg-gray-50 p-3">
            <p className="text-sm text-gray-600">Expired</p>
            <p className="text-lg font-semibold">
              {room.isExpired ? "Yes ⚠️" : "No ✅"}
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
