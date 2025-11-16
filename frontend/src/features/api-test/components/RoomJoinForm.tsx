import { useEffect, useRef } from "react";

type Props = {
  onJoinRoom: (userId: string, roomNumber: number) => void;
  loading: boolean;
  error: string | null;
  success: boolean;
  autoFillRoomNumber?: number | null;
};

/**
 * 部屋参加フォーム (Presentational)
 */
export function RoomJoinForm({
  onJoinRoom,
  loading,
  error,
  success,
  autoFillRoomNumber,
}: Props) {
  const roomNumberInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (autoFillRoomNumber && roomNumberInputRef.current) {
      roomNumberInputRef.current.value = String(autoFillRoomNumber);
    }
  }, [autoFillRoomNumber]);

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const userId = formData.get("userId") as string;
    const roomNumber = Number(formData.get("roomNumber"));
    onJoinRoom(userId, roomNumber);
  };

  return (
    <div className="rounded-lg border p-4">
      <h2 className="mb-4 text-xl font-bold">2. Join Room</h2>

      <form onSubmit={handleSubmit} className="space-y-3">
        <div>
          <label className="block text-sm font-medium">User ID:</label>
          <input
            name="userId"
            type="text"
            defaultValue="550e8400-e29b-41d4-a716-446655440002"
            className="w-full rounded border px-3 py-2"
            disabled={loading}
          />
        </div>

        <div>
          <label className="block text-sm font-medium">Room Number:</label>
          <input
            ref={roomNumberInputRef}
            name="roomNumber"
            type="number"
            defaultValue="1"
            className="w-full rounded border px-3 py-2"
            disabled={loading}
          />
        </div>

        <button
          type="submit"
          disabled={loading}
          className="rounded bg-green-500 px-4 py-2 text-white hover:bg-green-600 disabled:bg-gray-300"
        >
          {loading ? "Joining..." : "Join Room"}
        </button>
      </form>

      {error && (
        <div className="mt-3 rounded bg-red-100 p-3 text-red-700">
          Error: {error}
        </div>
      )}

      {success && (
        <div className="mt-3 rounded bg-green-100 p-3 text-green-700">
          ✅ Joined successfully!
        </div>
      )}
    </div>
  );
}
