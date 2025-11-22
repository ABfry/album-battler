import type { GetBattleResultResponse, UserInfo } from "@/src/lib/api/types";
import { Button } from "@/src/components/ui/button";
import { motion } from "framer-motion";
import Image from "next/image";

type ResultViewProps = {
  battleResult: GetBattleResultResponse | null;
  players: UserInfo[];
  theme: string | null;
  isLoading: boolean;
  onBackToTitle: () => void;
  onShare: () => void;
};

/**
 * 結果画面のプレゼンテーショナルコンポーネント (Presentational)
 * 表示のみを担当
 */
export function ResultView({
  battleResult,
  players,
  theme,
  isLoading,
  onBackToTitle,
  onShare,
}: ResultViewProps) {
  if (isLoading) {
    return (
      <main className="flex min-h-screen items-center justify-center">
        <div className="text-2xl font-bold text-slate-700">読み込み中...</div>
      </main>
    );
  }

  if (!battleResult) {
    return (
      <main className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <div className="mb-4 text-2xl font-bold text-red-600">
            結果の取得に失敗しました
          </div>
          <Button onClick={onBackToTitle}>タイトルに戻る</Button>
        </div>
      </main>
    );
  }

  const winnerPlayer = players.find(
    (p) => p.id === battleResult.winner_user_id
  );

  const getRankStyle = (rank: number) => {
    switch (rank) {
      case 1:
        return "bg-yellow-100 border-yellow-400 text-yellow-900";
      case 2:
        return "bg-gray-200 border-gray-400 text-gray-800";
      case 3:
        return "bg-amber-100 border-amber-400 text-amber-900";
      default:
        return "bg-white border-[#3551b8] text-slate-800";
    }
  };

  return (
    <main className="flex min-h-screen items-center justify-center p-4">
      {/* 戻るボタン */}
      <button
        onClick={onBackToTitle}
        className="absolute top-4 left-4 flex h-10 w-10 items-center justify-center rounded-md bg-white text-xl shadow hover:bg-gray-50"
      >
        ◀
      </button>

      {/* 結果コンテナ */}
      <div className="w-full max-w-4xl space-y-8">
        {/* お題表示 */}
        {theme && (
          <div className="text-center">
            <h2 className="mb-2 text-2xl font-bold text-slate-700">お題</h2>
            <p className="text-3xl font-black text-[#b57c39]">{theme}</p>
          </div>
        )}

        {/* 勝者発表 */}
        <div className="space-y-4 text-center">
          <div className="flex justify-center">
            <Image
              src="/winner-text.png"
              alt="Winner"
              width={260}
              height={100}
              priority
            />
          </div>
          <div className="text-5xl font-bold text-slate-800">
            {winnerPlayer?.name || "Unknown Player"}
          </div>
        </div>

        {/* 詳細結果カード */}
        <div className="space-y-4">
          {battleResult.results.map((result) => {
            const player = players.find((p) => p.id === result.user_id);
            const isWinner = result.user_id === battleResult.winner_user_id;
            const rankStyle = getRankStyle(result.rank);

            return (
              <div
                key={result.user_id}
                className={`rounded-xl border p-6 shadow-lg transition-all ${rankStyle} ${
                  isWinner ? "ring-4 ring-yellow-200" : ""
                }`}
              >
                {/* ヘッダー：順位・名前・スコア */}
                <div className="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                  {/* ユーザー情報 */}
                  <div className="flex items-center gap-3">
                    {/* 順位 */}
                    <span className={`text-2xl font-bold sm:text-3xl`}>
                      {result.rank}位
                    </span>
                    {/* ユーザーアイコン */}
                    <div className="relative h-10 w-10 shrink-0 overflow-hidden rounded-full border-2 border-gray-300 sm:h-12 sm:w-12">
                      {player?.icon_url ? (
                        <Image
                          src={player.icon_url}
                          alt={`${player.name}のアイコン`}
                          fill
                          className="object-cover"
                        />
                      ) : (
                        <div className="flex h-full w-full items-center justify-center bg-gray-200">
                          <span className="text-lg sm:text-xl">👤</span>
                        </div>
                      )}
                    </div>
                    {/* 名前と王冠 */}
                    <div className="flex items-center gap-2">
                      <span className="text-lg font-black text-slate-800 sm:text-2xl">
                        {player?.name || "Unknown"}
                      </span>
                      {isWinner && (
                        <div className="relative h-8 w-8 sm:h-10 sm:w-10">
                          <div className="absolute inset-0 flex items-center justify-center">
                            <motion.div
                              className="h-[120%] w-[120%]"
                              initial={{ rotate: 0 }}
                              animate={{ rotate: 360 }}
                              transition={{
                                duration: 3,
                                repeat: Infinity,
                                ease: "linear",
                              }}
                            >
                              <Image
                                src="/back-light.png"
                                alt="back light"
                                fill
                                unoptimized
                              />
                            </motion.div>
                          </div>
                          <div className="relative flex h-full w-full items-center justify-center">
                            <Image
                              src="/crown-icon.png"
                              alt="crown"
                              fill
                              unoptimized
                            />
                          </div>
                        </div>
                      )}
                    </div>
                  </div>
                  {/* スコア */}
                  <div className="flex items-baseline gap-2 sm:block sm:text-right">
                    <div className="text-2xl font-black text-[#b57c39] sm:text-4xl">
                      {Math.round(result.final_score)}点
                    </div>
                    <div className="text-xs text-gray-600 sm:text-sm">
                      AI: {Math.round(result.ai_score)}点 / 拍手:{" "}
                      {result.user_score}点
                    </div>
                  </div>
                </div>

                {/* 画像 */}
                <div className="relative mb-4 h-64 w-full overflow-hidden rounded-lg border border-gray-300">
                  <Image
                    src={result.image_url}
                    alt={`${player?.name}の画像`}
                    fill
                    className="object-contain"
                  />
                </div>

                {/* AI評価コメント */}
                {result.ai_explanation && (
                  <div className="rounded-lg bg-gray-50 p-4">
                    <h3 className="mb-2 font-bold text-slate-700">AI評価</h3>
                    <p className="text-sm leading-relaxed text-slate-600">
                      {result.ai_explanation}
                    </p>
                  </div>
                )}
              </div>
            );
          })}
        </div>

        {/* ボタン群 */}
        <div className="flex flex-col items-center gap-3 pb-8">
          <Button onClick={onShare} variant="secondary">
            この結果を共有
          </Button>
          <Button onClick={onBackToTitle}>タイトルに戻る</Button>
        </div>
      </div>
    </main>
  );
}
