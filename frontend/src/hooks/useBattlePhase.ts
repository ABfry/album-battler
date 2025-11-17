import { useState, useCallback } from "react";

/**
 * バトルのフェーズ定義
 */
export type BattlePhase =
  | "waiting" // ゲーム開始待ち
  | "selecting" // 画像選択中
  | "clap_time" // 拍手タイム
  | "result" // 結果発表
  | "finished"; // 終了

/**
 * 各フェーズの設定
 */
export type PhaseConfig = {
  duration: number; // このフェーズの制限時間（秒）
  canSelectImage: boolean; // 画像選択可能か
  canClap: boolean; // 拍手可能か
};

/**
 * フェーズごとの設定値
 */
export const PHASE_CONFIGS: Record<BattlePhase, PhaseConfig> = {
  waiting: { duration: 0, canSelectImage: false, canClap: false },
  selecting: { duration: 30, canSelectImage: true, canClap: false },
  clap_time: { duration: 60, canSelectImage: false, canClap: true },
  result: { duration: 10, canSelectImage: false, canClap: false },
  finished: { duration: 0, canSelectImage: false, canClap: false },
};

/**
 * バトルフェーズ管理フック
 * WebSocketイベントに基づくフェーズ遷移を管理
 */
export function useBattlePhase(initialPhase: BattlePhase = "waiting") {
  const [phase, setPhase] = useState<BattlePhase>(initialPhase);

  const config = PHASE_CONFIGS[phase];

  // フェーズ遷移
  const transitionTo = useCallback(
    (newPhase: BattlePhase) => {
      console.log(`[useBattlePhase] Phase transition: ${phase} -> ${newPhase}`);
      setPhase(newPhase);
    },
    [phase]
  );

  return {
    phase,
    config,
    transitionTo,
    canSelectImage: config.canSelectImage,
    canClap: config.canClap,
  };
}
