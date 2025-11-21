"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useResult } from "@/src/hooks/useResult";
import { useRoomInfo } from "@/src/hooks/useRoomInfo";
import { useBattleInfo } from "@/src/hooks/useBattleInfo";
import type { GetBattleResultResponse } from "@/src/lib/api/types";
import { ResultView } from "../components/ResultView";

type ResultPageProps = {
  battleId: string;
};

/**
 * 結果画面のコンテナコンポーネント (Container)
 * ロジック・状態管理を担当
 */
export function ResultPage({ battleId }: ResultPageProps) {
  const router = useRouter();
  const { getResult } = useResult();
  const [battleResult, setBattleResult] =
    useState<GetBattleResultResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // バトル情報取得（お題表示用）
  const { battle } = useBattleInfo(battleId);

  // ルーム情報取得（プレイヤー情報用）
  const { room } = useRoomInfo(battle?.roomId || null);

  useEffect(() => {
    const fetchResult = async () => {
      setIsLoading(true);
      const result = await getResult(battleId);
      if (result) {
        setBattleResult(result);
      }
      setIsLoading(false);
    };

    fetchResult();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [battleId]);

  const handleBackToTitle = () => {
    router.push("/title");
  };

  const players = room?.users || [];

  return (
    <ResultView
      battleResult={battleResult}
      players={players}
      theme={battle?.theme || null}
      isLoading={isLoading}
      onBackToTitle={handleBackToTitle}
    />
  );
}
