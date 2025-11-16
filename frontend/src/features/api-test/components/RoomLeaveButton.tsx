type RoomLeaveButtonProps = {
  roomId: string | null;
  userId: string;
  onLeaveRoom: (roomId: string, userId: string) => Promise<void>;
  loading: boolean;
  error: string | null;
  success: boolean;
};

export function RoomLeaveButton({
  roomId,
  userId,
  onLeaveRoom,
  loading,
  error,
  success,
}: RoomLeaveButtonProps) {
  const handleLeave = async () => {
    if (!roomId) return;
    await onLeaveRoom(roomId, userId);
  };

  return (
    <div className="rounded-lg border bg-white p-4">
      <h2 className="mb-4 text-xl font-bold">4. Leave Room</h2>

      <div className="space-y-4">
        <div>
          <label className="mb-1 block text-sm font-medium text-gray-700">
            Room ID:
          </label>
          <p className="text-sm text-gray-600">
            {roomId || "No room selected"}
          </p>
        </div>

        <div>
          <label className="mb-1 block text-sm font-medium text-gray-700">
            User ID:
          </label>
          <input
            type="text"
            value={userId}
            disabled
            className="w-full rounded-md border border-gray-300 bg-gray-50 px-3 py-2 text-sm"
          />
        </div>

        <button
          onClick={handleLeave}
          disabled={!roomId || loading}
          className="w-full rounded-md bg-red-600 px-4 py-2 text-white hover:bg-red-700 disabled:cursor-not-allowed disabled:bg-gray-400"
        >
          {loading ? "Leaving..." : "Leave Room"}
        </button>

        {error && (
          <div className="rounded-md bg-red-50 p-3 text-sm text-red-600">
            Error: {error}
          </div>
        )}

        {success && (
          <div className="rounded-md bg-green-50 p-3 text-sm text-green-600">
            ✅ Left successfully!
          </div>
        )}
      </div>
    </div>
  );
}
