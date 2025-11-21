import type { GetBattleResultResponse, UserInfo } from "@/src/lib/api/types";

type ResultDisplayProps = {
  battleResult: GetBattleResultResponse;
  players: UserInfo[];
};

/**
 * バトル結果表示コンポーネント (Presentational)
 */
export function ResultDisplay({ battleResult, players }: ResultDisplayProps) {
  const winnerPlayer = players.find(
    (p) => p.id === battleResult.winner_user_id
  );

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-linear-to-b from-yellow-50 to-orange-50 p-4">
      <div className="w-full max-w-4xl space-y-8">
        {/* 勝者発表 */}
        <div className="space-y-4 text-center">
          <h1 className="bg-linear-to-r from-yellow-400 to-orange-500 bg-clip-text text-6xl font-black text-transparent">
            WINNER
          </h1>
          <div className="text-5xl font-bold text-slate-800">
            {winnerPlayer?.name || "Unknown Player"}
          </div>
        </div>

        {/* スコアボード */}
        <div className="rounded-2xl bg-white p-8 shadow-2xl">
          <h2 className="mb-6 text-center text-3xl font-bold text-slate-700">
            スコア
          </h2>
          <div className="space-y-3">
            {battleResult.results.map((result) => {
              const player = players.find((p) => p.id === result.user_id);
              const isWinner = result.user_id === battleResult.winner_user_id;

              return (
                <div
                  key={result.user_id}
                  className={`flex items-center justify-between rounded-lg p-5 transition-all ${
                    isWinner
                      ? "scale-105 bg-gradient-to-r from-yellow-300 via-yellow-200 to-yellow-300 shadow-lg"
                      : "bg-gray-50 hover:bg-gray-100"
                  }`}
                >
                  <div className="flex items-center gap-4">
                    <span
                      className={`text-2xl font-bold ${
                        isWinner ? "text-yellow-700" : "text-gray-500"
                      }`}
                    >
                      {result.rank}
                    </span>
                    <span
                      className={`text-xl font-semibold ${
                        isWinner ? "text-yellow-900" : "text-slate-700"
                      }`}
                    >
                      {player?.name || "Unknown"}
                    </span>
                    {isWinner && <span className="text-2xl">👑</span>}
                  </div>
                  <span
                    className={`text-3xl font-bold ${
                      isWinner ? "text-yellow-700" : "text-slate-600"
                    }`}
                  >
                    {Math.round(result.final_score)}点
                  </span>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}
