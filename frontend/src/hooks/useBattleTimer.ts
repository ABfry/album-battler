import { useState, useEffect, useCallback, useRef } from "react";

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

  const { onWarning, onTimeUp } = options;

  // コールバックをRefで保持（依存配列から除外するため）
  const onWarningRef = useRef(onWarning);
  const onTimeUpRef = useRef(onTimeUp);

  useEffect(() => {
    onWarningRef.current = onWarning;
    onTimeUpRef.current = onTimeUp;
  }, [onWarning, onTimeUp]);

  // カウントダウン処理
  useEffect(() => {
    if (!isRunning) return;

    const timer = setInterval(() => {
      setTimeLeft((prev) => {
        // すでに0なら処理しない
        if (prev <= 0) return 0;

        const next = Math.max(0, prev - 1);

        // 警告チェック（残り10秒以下）
        if (next <= 10 && next > 0) {
          setIsWarning(true);
          onWarningRef.current?.(next);
        }

        // タイムアップ
        if (next === 0) {
          onTimeUpRef.current?.();
          setIsRunning(false);
        }

        return next;
      });
    }, 1000);

    return () => clearInterval(timer);
  }, [isRunning]);

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
