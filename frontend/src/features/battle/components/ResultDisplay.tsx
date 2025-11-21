import type { GetBattleResultResponse, UserInfo } from "@/src/lib/api/types";
import Image from "next/image";

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
    <div className="flex min-h-screen flex-col items-center justify-center p-4">
      <div className="w-full max-w-md space-y-8">
        {/* 勝者発表 */}
        <div className="space-y-4 text-center">
          <h1 className="text-5xl font-black tracking-widest text-[#b57c39]">
            WINNER
          </h1>
          <div className="text-4xl font-black text-slate-800">
            {winnerPlayer?.name || "Unknown Player"}
          </div>
        </div>

        {/* スコアボード */}
        <div className="rounded-xl border border-[#3551b8] bg-white p-6 shadow-[0_8px_0_rgba(0,0,0,0.15)]">
          <h2 className="mb-4 text-center text-2xl font-black text-slate-700">
            スコア
          </h2>
          <div className="space-y-3">
            {battleResult.results.map((result) => {
              const player = players.find((p) => p.id === result.user_id);
              const isWinner = result.user_id === battleResult.winner_user_id;

              return (
                <div
                  key={result.user_id}
                  className={`flex items-center justify-between rounded-lg p-4 ${
                    isWinner
                      ? "border-2 border-yellow-400 bg-yellow-50"
                      : "bg-gray-50"
                  }`}
                >
                  <div className="flex items-center gap-3">
                    <span
                      className={`text-2xl font-black ${
                        isWinner ? "text-yellow-600" : "text-gray-500"
                      }`}
                    >
                      {result.rank}
                    </span>
                    {/* ユーザーアイコン */}
                    <div className="relative h-10 w-10 overflow-hidden rounded-full border-2 border-gray-300">
                      {player?.icon_url ? (
                        <Image
                          src={player.icon_url}
                          alt={`${player.name}のアイコン`}
                          fill
                          className="object-cover"
                        />
                      ) : (
                        <div className="flex h-full w-full items-center justify-center bg-gray-200">
                          <span className="text-lg">👤</span>
                        </div>
                      )}
                    </div>
                    <span
                      className={`text-lg font-black ${
                        isWinner ? "text-slate-800" : "text-slate-700"
                      }`}
                    >
                      {player?.name || "Unknown"}
                    </span>
                    {isWinner && <span className="text-2xl">👑</span>}
                  </div>
                  <span
                    className={`text-2xl font-black ${
                      isWinner ? "text-[#b57c39]" : "text-slate-600"
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
