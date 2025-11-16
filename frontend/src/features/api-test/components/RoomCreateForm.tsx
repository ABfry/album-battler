type Props = {
  onCreateRoom: (userId: string) => void;
  loading: boolean;
  error: string | null;
  result: { room_id: string; room_number: number } | null;
};

/**
 * 部屋作成フォーム (Presentational)
 */
export function RoomCreateForm({
  onCreateRoom,
  loading,
  error,
  result,
}: Props) {
  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const userId = formData.get("userId") as string;
    onCreateRoom(userId);
  };

  return (
    <div className="rounded-lg border bg-white p-4">
      <h2 className="mb-4 text-xl font-bold">1. Create Room</h2>

      <form onSubmit={handleSubmit} className="space-y-3">
        <div>
          <label className="block text-sm font-medium">User ID:</label>
          <input
            name="userId"
            type="text"
            defaultValue="550e8400-e29b-41d4-a716-446655440001"
            className="w-full rounded border px-3 py-2"
            disabled={loading}
          />
        </div>

        <button
          type="submit"
          disabled={loading}
          className="rounded bg-blue-500 px-4 py-2 text-white hover:bg-blue-600 disabled:bg-gray-300"
        >
          {loading ? "Creating..." : "Create Room"}
        </button>
      </form>

      {error && (
        <div className="mt-3 rounded bg-red-100 p-3 text-red-700">
          Error: {error}
        </div>
      )}

      {result && (
        <div className="mt-3 rounded bg-green-100 p-3 text-green-700">
          <p>✅ Room Created!</p>
          <p>Room ID: {result.room_id}</p>
          <p>Room Number: {result.room_number}</p>
        </div>
      )}
    </div>
  );
}
