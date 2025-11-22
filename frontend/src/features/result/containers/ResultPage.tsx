"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useResult } from "@/src/hooks/useResult";
import { useRoomInfo } from "@/src/hooks/useRoomInfo";
import { useBattleInfo } from "@/src/hooks/useBattleInfo";
import type { GetBattleResultResponse } from "@/src/lib/api/types";
import { ResultView } from "../components/ResultView";
import { toast } from "sonner";

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
    router.push("/");
  };

  const handleShare = async () => {
    const url = window.location.href;
    const shareData = {
      title: "アルバムバトラー - 対戦結果",
      text: battle?.theme
        ? `お題「${battle.theme}」のバトル結果をチェック！`
        : "バトル結果をチェック！",
      url,
    };

    // Web Share APIが利用可能な場合はそれを使用
    if (navigator.share) {
      try {
        await navigator.share(shareData);
      } catch (err) {
        // ユーザーがキャンセルした場合など
        if ((err as Error).name !== "AbortError") {
          // フォールバック: クリップボードにコピー
          await copyToClipboard(url);
        }
      }
    } else {
      // フォールバック: クリップボードにコピー
      await copyToClipboard(url);
    }
  };

  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
      toast.success("URLをコピーしました");
    } catch {
      toast.error("コピーに失敗しました");
    }
  };

  const players = room?.users || [];

  return (
    <ResultView
      battleResult={battleResult}
      players={players}
      theme={battle?.theme || null}
      isLoading={isLoading}
      onBackToTitle={handleBackToTitle}
      onShare={handleShare}
    />
  );
}
