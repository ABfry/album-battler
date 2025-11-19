import Image from "next/image";
import type { BattleImage, UserInfo } from "@/src/lib/api/types";

type PlayerListProps = {
  players: UserInfo[] | undefined;
  images: BattleImage[];
};

/**
 * プレイヤー一覧表示コンポーネント (Presentational)
 */
export function PlayerList({ players, images }: PlayerListProps) {
  if (!players || players.length === 0) {
    return null;
  }

  // ユーザーが画像を提出済みかチェック
  const hasSubmitted = (userId: string) => {
    return images?.some((img) => img.userId === userId) ?? false;
  };

  return (
    <div className="w-full">
      {/* アイコン */}
      <div className="mb-4 flex items-center justify-center gap-4 md:gap-6">
        {players.map((player) => (
          <div key={player.id} className="flex flex-col items-center gap-2">
            {/* ステータスアイコン */}
            <div className="text-4xl">
              {hasSubmitted(player.id) ? "✅" : "🧐"}
            </div>

            {/* アイコン画像 */}
            <div className="h-16 w-16 overflow-hidden rounded-full border-4 border-slate-300 bg-slate-200 shadow-lg md:h-20 md:w-20">
              <Image
                src={player.icon_url}
                alt={player.name}
                width={80}
                height={80}
                className="h-full w-full object-cover"
                unoptimized
              />
            </div>
          </div>
        ))}
      </div>

      {/* ユーザー名 */}
      <div className="flex items-center justify-center gap-4 md:gap-6">
        {players.map((player) => (
          <div key={player.id} className="w-16 text-center md:w-20">
            <p className="truncate text-sm font-bold text-slate-700 md:text-base">
              {player.name}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}
