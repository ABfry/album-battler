import { useState, useEffect, useCallback } from "react";

type UseBattleTimerOptions = {
  initialTime?: number;
  autoStart?: boolean;
  onWarning?: (secondsLeft: number) => void;
  onTimeUp?: () => void;
};

/**
 * バトルタイマー管理フック
 * カウントダウン表示と警告状態の管理を担当
 * フェーズ管理はしない（WebSocketに任せる）
 */
export function useBattleTimer(options: UseBattleTimerOptions = {}) {
  const [timeLeft, setTimeLeft] = useState(options.initialTime ?? 30);
  const [isRunning, setIsRunning] = useState(options.autoStart ?? false);
  const [isWarning, setIsWarning] = useState(false);

  // カウントダウン処理
  useEffect(() => {
    if (!isRunning || timeLeft <= 0) return;

    const timer = setInterval(() => {
      setTimeLeft((prev) => {
        const next = Math.max(0, prev - 1);

        // 警告チェック（残り10秒以下）
        if (next <= 10 && !isWarning) {
          setIsWarning(true);
          options.onWarning?.(next);
        }

        // タイムアップ
        if (next === 0) {
          options.onTimeUp?.();
          setIsRunning(false);
        }

        return next;
      });
    }, 1000);

    return () => clearInterval(timer);
  }, [isRunning, timeLeft, isWarning, options]);

  const startTimer = useCallback(() => {
    setIsRunning(true);
  }, []);

  const pauseTimer = useCallback(() => {
    setIsRunning(false);
  }, []);

  const resetTimer = useCallback((time: number) => {
    setTimeLeft(time);
    setIsWarning(false);
  }, []);

  const setTime = useCallback((time: number) => {
    setTimeLeft(time);
  }, []);

  return {
    timeLeft,
    isWarning,
    isRunning,
    startTimer,
    pauseTimer,
    resetTimer,
    setTime,
  };
}
