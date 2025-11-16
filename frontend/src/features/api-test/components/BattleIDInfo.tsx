type Props = {
  battleID: string | null;
};

/**
 * バトルID表示 (Presentational)
 * start_gameイベント受信時に自動取得されたBattle IDを表示
 */
export function BattleIDInfo({ battleID }: Props) {
  return (
    <div className="rounded-lg border p-4 shadow-sm">
      <h2 className="mb-4 text-xl font-bold">6. Get Battle ID</h2>

      <div className="mb-3 space-y-2 rounded bg-gray-50 p-3">
        <p className="text-sm text-gray-600">
          💡 start_gameイベント受信時に自動取得されます
        </p>
      </div>

      {battleID ? (
        <div className="rounded bg-green-50 p-3">
          <p className="text-sm font-semibold text-green-700">Battle ID:</p>
          <p className="font-mono text-xs break-all text-green-600">
            {battleID}
          </p>
        </div>
      ) : (
        <div className="rounded bg-gray-100 p-3">
          <p className="text-sm text-gray-500">ゲーム開始後に表示されます</p>
        </div>
      )}
    </div>
  );
}
