import type { GetBattleResultResponse, UserInfo } from "@/src/lib/api/types";
import { Button } from "@/src/components/ui/button";
import Image from "next/image";

type ResultViewProps = {
  battleResult: GetBattleResultResponse | null;
  players: UserInfo[];
  theme: string | null;
  isLoading: boolean;
  onBackToTitle: () => void;
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
          <h1 className="bg-linear-to-r from-yellow-400 to-orange-500 bg-clip-text text-6xl font-black text-transparent">
            WINNER
          </h1>
          <div className="text-5xl font-bold text-slate-800">
            {winnerPlayer?.name || "Unknown Player"}
          </div>
        </div>

        {/* 詳細結果カード */}
        <div className="space-y-4">
          {battleResult.results.map((result) => {
            const player = players.find((p) => p.id === result.user_id);
            const isWinner = result.user_id === battleResult.winner_user_id;

            return (
              <div
                key={result.user_id}
                className={`rounded-xl border bg-white p-6 shadow-lg transition-all ${
                  isWinner
                    ? "border-yellow-400 ring-4 ring-yellow-200"
                    : "border-[#3551b8]"
                }`}
              >
                {/* ヘッダー：順位・名前・スコア */}
                <div className="mb-4 flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <span
                      className={`text-3xl font-bold ${
                        isWinner ? "text-yellow-600" : "text-gray-500"
                      }`}
                    >
                      {result.rank}位
                    </span>
                    {/* ユーザーアイコン */}
                    <div className="relative h-12 w-12 overflow-hidden rounded-full border-2 border-gray-300">
                      {player?.icon_url ? (
                        <Image
                          src={player.icon_url}
                          alt={`${player.name}のアイコン`}
                          fill
                          className="object-cover"
                        />
                      ) : (
                        <div className="flex h-full w-full items-center justify-center bg-gray-200">
                          <span className="text-xl">👤</span>
                        </div>
                      )}
                    </div>
                    <span className="text-2xl font-black text-slate-800">
                      {player?.name || "Unknown"}
                    </span>
                    {isWinner && <span className="text-3xl">👑</span>}
                  </div>
                  <div className="text-right">
                    <div className="text-4xl font-black text-[#b57c39]">
                      {Math.round(result.final_score)}点
                    </div>
                    <div className="text-sm text-gray-600">
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

        {/* タイトルに戻るボタン */}
        <div className="flex justify-center pb-8">
          <Button onClick={onBackToTitle}>タイトルに戻る</Button>
        </div>
      </div>
    </main>
  );
}
