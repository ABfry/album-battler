import { useState, useCallback, useRef, useEffect } from "react";

/**
 * バトルのフェーズ定義
 */
export type BattlePhase =
  | "waiting" // ゲーム開始待ち
  | "selecting" // 画像選択中
  | "clap_time_1"
  | "clap_time_2"
  | "clap_time_3"
  | "clap_time_4"
  | "clap_time_5"
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
  selecting: { duration: 60, canSelectImage: true, canClap: false },
  clap_time_1: { duration: 5, canSelectImage: false, canClap: true },
  clap_time_2: { duration: 5, canSelectImage: false, canClap: true },
  clap_time_3: { duration: 5, canSelectImage: false, canClap: true },
  clap_time_4: { duration: 5, canSelectImage: false, canClap: true },
  clap_time_5: { duration: 5, canSelectImage: false, canClap: true },
  result: { duration: 10, canSelectImage: false, canClap: false },
  finished: { duration: 0, canSelectImage: false, canClap: false },
};

// 拍手フェーズの配列
export const CLAP_PHASES: BattlePhase[] = [
  "clap_time_1",
  "clap_time_2",
  "clap_time_3",
  "clap_time_4",
  "clap_time_5",
];

export type PhaseEventHandlers = {
  onPhaseStart?: (phase: BattlePhase) => void | Promise<void>;
  onTimeUp?: (phase: BattlePhase) => void | Promise<void>;
  onPhaseEnd?: (phase: BattlePhase) => void | Promise<void>;
};

export type PhaseHandlersMap = Partial<Record<BattlePhase, PhaseEventHandlers>>;

type UseBattlePhaseOptions = {
  initialPhase?: BattlePhase;
  handlers?: PhaseHandlersMap;
};

/**
 * バトルフェーズ管理フック
 * WebSocketイベントに基づくフェーズ遷移を管理
 */
export function useBattlePhase(options: UseBattlePhaseOptions = {}) {
  const { initialPhase = "waiting", handlers } = options;

  const [phase, setPhase] = useState<BattlePhase>(initialPhase);
  const config = PHASE_CONFIGS[phase];

  // handlersをRefで保持（常に最新のハンドラーを使用）
  const handlersRef = useRef<PhaseHandlersMap>(handlers ?? {});

  // handlersが変更されたらRefを更新
  useEffect(() => {
    if (handlers) {
      handlersRef.current = handlers;
    }
  }, [handlers]);

  // 外部からハンドラーを設定するための関数
  const setHandlers = useCallback((nextHandlers: PhaseHandlersMap) => {
    handlersRef.current = nextHandlers;
  }, []);

  // フェーズ遷移
  const transitionTo = useCallback(async (newPhase: BattlePhase) => {
    setPhase((currentPhase) => {
      if (currentPhase === newPhase) {
        return currentPhase;
      }

      console.log(
        `[useBattlePhase] Phase transition: ${currentPhase} -> ${newPhase}`
      );

      // イベントハンドラーを非同期で実行（状態更新後に）
      Promise.resolve().then(async () => {
        // onPhaseEnd
        try {
          await handlersRef.current[currentPhase]?.onPhaseEnd?.(currentPhase);
        } catch (error) {
          console.error(`[useBattlePhase] onPhaseEnd error: ${error}`);
        }

        // onPhaseStart
        try {
          await handlersRef.current[newPhase]?.onPhaseStart?.(newPhase);
        } catch (error) {
          console.error(`[useBattlePhase] onPhaseStart error: ${error}`);
        }
      });

      return newPhase;
    });
  }, []);

  // タイムアップ
  const onTimeUp = useCallback(async () => {
    try {
      await handlersRef.current[phase]?.onTimeUp?.(phase);
    } catch (error) {
      console.error(`[useBattlePhase] onTimeUp error: ${error}`);
    }
  }, [phase]);

  return {
    phase,
    config,
    transitionTo,
    onTimeUp,
    setHandlers,
    canSelectImage: config.canSelectImage,
    canClap: config.canClap,
  };
}
