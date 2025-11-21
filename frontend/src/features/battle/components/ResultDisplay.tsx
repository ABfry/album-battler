import type { GetBattleResultResponse, UserInfo } from "@/src/lib/api/types";
import Image from "next/image";
import { motion } from "framer-motion";

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

  // 最大スコアを取得（棒グラフの高さ計算用）
  const maxScore = Math.max(...battleResult.results.map((r) => r.final_score));

  // プレイヤーごとの色定義（順番に割り当て）
  const playerColors = [
    { bar: "bg-gradient-to-t from-red-500 to-red-400", text: "text-red-600" },
    {
      bar: "bg-gradient-to-t from-blue-500 to-blue-400",
      text: "text-blue-600",
    },
    {
      bar: "bg-gradient-to-t from-yellow-500 to-yellow-400",
      text: "text-yellow-600",
    },
    {
      bar: "bg-gradient-to-t from-green-500 to-green-400",
      text: "text-green-600",
    },
    {
      bar: "bg-gradient-to-t from-orange-500 to-orange-400",
      text: "text-orange-600",
    },
  ];

  return (
    <div className="flex min-h-screen flex-col items-center justify-center p-4">
      <div className="relative flex h-screen w-full max-w-4xl flex-col items-center justify-center py-8">
        {/* 勝者発表（中央） */}
        <motion.div
          initial={{ opacity: 0, scale: 0.8 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{
            delay: battleResult.results.length * 0.2 + 1.2,
            duration: 0.8,
          }}
          className="mb-16 space-y-4 text-center"
        >
          <h1 className="text-5xl font-black tracking-widest text-[#b57c39]">
            WINNER
          </h1>
          <div className="text-4xl font-black text-slate-800">
            {winnerPlayer?.name || "Unknown Player"}
          </div>
        </motion.div>

        {/* プレイヤー一覧と棒グラフ（下部） */}
        <div className="w-full">
          {/* 棒グラフ */}
          <div className="mb-4 flex items-end justify-center gap-4 md:gap-6">
            {battleResult.results.map((result, index) => {
              const player = players.find((p) => p.id === result.user_id);
              const isWinner = result.user_id === battleResult.winner_user_id;
              // 棒グラフの高さ（最大200px）
              const barHeight = (result.final_score / maxScore) * 200;
              // プレイヤーの色を取得
              const color = playerColors[index % playerColors.length];

              return (
                <div
                  key={result.user_id}
                  className="flex flex-col items-center"
                >
                  {/* 点数表示 */}
                  <motion.div
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.2 + 1.0, duration: 0.5 }}
                    className={`mb-2 text-2xl font-black ${color.text}`}
                  >
                    {Math.round(result.final_score)}点
                  </motion.div>

                  {/* 棒グラフ */}
                  <motion.div
                    initial={{ height: 0 }}
                    animate={{ height: barHeight }}
                    transition={{
                      delay: index * 0.2,
                      duration: 1.0,
                      ease: "easeOut",
                    }}
                    className={`w-16 rounded-t-lg md:w-20 ${color.bar}`}
                    style={{ minHeight: "20px" }}
                  >
                    {/* 王冠（勝者のみ） */}
                    {isWinner && (
                      <motion.div
                        initial={{ opacity: 0, scale: 0 }}
                        animate={{ opacity: 1, scale: 1 }}
                        transition={{ delay: index * 0.2 + 0.8, duration: 0.3 }}
                        className="flex justify-center pt-2 text-3xl"
                      >
                        👑
                      </motion.div>
                    )}
                  </motion.div>

                  {/* アイコン画像 */}
                  <div className="mt-4 h-16 w-16 overflow-hidden rounded-full border-4 border-slate-300 bg-slate-200 shadow-lg md:h-20 md:w-20">
                    <Image
                      src={player?.icon_url || ""}
                      alt={player?.name || "Unknown"}
                      width={80}
                      height={80}
                      className="h-full w-full object-cover"
                      unoptimized
                    />
                  </div>

                  {/* ユーザー名 */}
                  <div className="mt-2 w-16 text-center md:w-20">
                    <p className="truncate text-sm font-bold text-slate-700 md:text-base">
                      {player?.name || "Unknown"}
                    </p>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}
