"use client";

import { useState, useEffect, useCallback, useMemo, useRef } from "react";
import { getUserIdClient } from "@/src/lib/auth/getUserIdClient";
import { useBattleInfo } from "@/src/hooks/useBattleInfo";
import { useRoomInfo } from "@/src/hooks/useRoomInfo";
import { useResult } from "@/src/hooks/useResult";
import {
  useBattlePhase,
  PHASE_CONFIGS,
  PhaseHandlersMap,
  CLAP_PHASES,
  BattlePhase,
} from "@/src/hooks/useBattlePhase";
import { useBattleTimer } from "@/src/hooks/useBattleTimer";
import { useBattleWebSocket } from "@/src/hooks/useBattleWebSocket";
import { useImageSelection } from "@/src/hooks/useImageSelection";
import { useWebSocketClap } from "@/src/lib/websocket/hooks/useWebSocketClap";
import { useWebSocket } from "@/src/lib/websocket/contexts/WebSocketContext";
import { useGameStateRestore } from "@/src/hooks/useGameStateRestore";
import { Battle } from "../components/Battle";
import type { ClapEffect } from "../components/Battle";
import type { GetBattleResultResponse } from "@/src/lib/api/types";

type BattlePageProps = {
  battleID: string;
};

/**
 * バトル画面のコンテナコンポーネント (Container)
 * ロジック・状態管理を担当
 */
export function BattlePage({ battleID }: BattlePageProps) {
  // CookieからユーザーIDを取得（フォールバックは固定値）
  const [userId] = useState(() => {
    const userId = getUserIdClient() || "";
    if (userId === "") {
      console.error("ユーザーIDがありません");
      return "";
    }
    return userId;
  });

  const [showErrorDialog, setShowErrorDialog] = useState(false);
  // 拍手フェーズで表示する画像（undefined: 通常モード, null: 未提出, string: 画像URL）
  const [displayedImage, setDisplayedImage] = useState<
    string | null | undefined
  >(undefined);
  const [battleResult, setBattleResult] =
    useState<GetBattleResultResponse | null>(null);
  const [remoteClapEffects, setRemoteClapEffects] = useState<ClapEffect[]>([]);

  // 1. バトル情報取得
  const {
    battle,
    images,
    loading: battleInfoLoading,
    error: battleInfoError,
    refetch: refetchBattle,
    refetchImages,
  } = useBattleInfo(battleID);

  // 2. 部屋情報取得
  const { room, refetch: refetchRoom } = useRoomInfo(battle?.roomId || null);

  const players = useMemo(() => room?.users ?? [], [room?.users]);

  // WebSocket再接続用: roomIdとbattleIdを保存
  const { setCurrentRoomId, setCurrentBattleId } = useWebSocket();

  useEffect(() => {
    if (battle?.roomId) {
      setCurrentRoomId(battle.roomId);
    }
  }, [battle?.roomId, setCurrentRoomId]);

  useEffect(() => {
    if (battleID) {
      setCurrentBattleId(battleID);
    }
  }, [battleID, setCurrentBattleId]);

  // ゲーム状態復元（WebSocket再接続時）
  const { restoredState } = useGameStateRestore();

  // 3. 結果取得
  const { getResult } = useResult();

  // 3. フェーズ管理
  const battlePhase = useBattlePhase({
    initialPhase: "waiting",
  });

  // 4. タイマー
  const timer = useBattleTimer({
    initialTime: battlePhase.config.duration,
    autoStart: false,
    onWarning: useCallback((secondsLeft: number) => {
      console.log(`[BattlePage] Warning: ${secondsLeft} seconds left`);
    }, []),
    onTimeUp: battlePhase.onTimeUp,
  });

  // 状態復元処理（WebSocket再接続時）
  useEffect(() => {
    if (!restoredState) return;

    console.log("[BattlePage] Restoring game state:", restoredState);

    const { current_phase, clap_current_user_index, selecting_started_at } =
      restoredState;

    // フェーズ復元
    if (current_phase === "selecting") {
      battlePhase.transitionTo("selecting");

      // タイマー復元: 開始時刻から経過時間を計算
      if (selecting_started_at) {
        const startTime = new Date(selecting_started_at).getTime();
        const now = Date.now();
        const elapsed = Math.floor((now - startTime) / 1000);
        const remainingTime = Math.max(60 - elapsed, 0);

        timer.resetTimer(remainingTime);
        if (remainingTime > 0) {
          timer.startTimer();
        }
      }
    } else if (current_phase === "clap_time") {
      // 拍手フェーズ復元: clap_current_user_indexに基づいてフェーズを決定
      const clapIndex = clap_current_user_index ?? 0;
      const clapPhase =
        CLAP_PHASES[clapIndex] || ("clap_time_1" as BattlePhase);

      battlePhase.transitionTo(clapPhase);

      // タイマーは各拍手フェーズのonPhaseStartで再設定されるため、ここでは明示的な復元不要
    } else if (current_phase === "result") {
      battlePhase.transitionTo("result");
    }

    console.log(`[BattlePage] Game state restored to phase: ${current_phase}`);
  }, [restoredState, battlePhase, timer]);

  const playersRef = useRef(players);
  const imagesRef = useRef(images);
  const battleRef = useRef(battle);

  useEffect(() => {
    playersRef.current = players;
  }, [players]);

  useEffect(() => {
    imagesRef.current = images;
  }, [images]);

  useEffect(() => {
    battleRef.current = battle;
  }, [battle]);

  useEffect(() => {
    // timerとbatttlePhaseがハンドラ作成時点で存在していないため，
    // 循環参照を避けるためにuseEffect内で実行
    const handlers: PhaseHandlersMap = {
      waiting: {
        onPhaseStart: () => console.log("ゲーム開始待ち"),
      },
      selecting: {
        onPhaseStart: () => {
          console.log("画像選択開始");
          const timeLimit = battleRef.current?.battleTimeLimitSeconds ?? 60;
          timer.resetTimer(timeLimit);
          timer.startTimer();
        },
        onTimeUp: () => {
          console.log("選択時間終了");
        },
        onPhaseEnd: () => {
          console.log("選択終了");
        },
      },
      result: {
        onPhaseStart: () => {
          console.log("結果発表");
          // 拍手フェーズ終了、displayedImageをクリア
          setDisplayedImage(undefined);
          // WebSocketイベント駆動で結果取得するため、ここでは何もしない
        },
      },
    };

    CLAP_PHASES.forEach((clapPhase, index) => {
      const nextPhase = CLAP_PHASES[index + 1] ?? "result";

      handlers[clapPhase] = {
        onPhaseStart: () => {
          const currentPlayers = playersRef.current;
          const currentImages = imagesRef.current;
          const player = currentPlayers[index];

          if (!player) {
            console.log(`${clapPhase}: プレイヤーがいないのでスキップ`);
            battlePhase.transitionTo(nextPhase as BattlePhase);
            return;
          }

          // プレイヤーの画像を表示
          const playerImage = currentImages.find(
            (img) => img.userId === player.id
          );
          if (playerImage) {
            setDisplayedImage(playerImage.imageUrl);
            console.log(`${player.name}の画像を表示: ${playerImage.imageUrl}`);
          } else {
            setDisplayedImage(null);
            console.log(`${player.name}の画像が見つかりません`);
          }

          console.log(
            `${player.name}の拍手開始 (${index + 1}/${currentPlayers.length})`
          );
          const duration = PHASE_CONFIGS[clapPhase].duration;
          timer.resetTimer(duration);
          timer.startTimer();
        },
        onTimeUp: () => {
          const currentPlayers = playersRef.current;
          const player = currentPlayers[index];
          console.log(`${player?.name || "Player"}の拍手終了`);
          battlePhase.transitionTo(nextPhase as BattlePhase);
        },
        onPhaseEnd: () => {
          const currentPlayers = playersRef.current;
          const player = currentPlayers[index];
          console.log(`${player?.name || "Player"}の拍手フェーズ完了`);
        },
      };
    });

    battlePhase.setHandlers(handlers);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const phaseMessage = useMemo(() => {
    if (battlePhase.phase === "selecting") {
      return timer.timeLeft === 0 ? "タイムアップ！" : "画像を探せ！";
    }

    const clapIndex = CLAP_PHASES.indexOf(battlePhase.phase);
    if (clapIndex !== -1) {
      const player = players[clapIndex];
      return player ? `${player.name}の画像` : "拍手タイム";
    }

    if (battlePhase.phase === "result") {
      return "結果発表";
    }

    return "";
  }, [battlePhase.phase, players, timer.timeLeft]);

  const [toastMessage, setToastMessage] = useState("");
  const [showToast, setShowToast] = useState(false);

  useEffect(() => {
    if (!phaseMessage) {
      setShowToast(false);
      return;
    }

    setToastMessage(phaseMessage);
    setShowToast(true);

    const timeoutId = window.setTimeout(() => {
      setShowToast(false);
    }, 2000);

    return () => window.clearTimeout(timeoutId);
  }, [phaseMessage]);

  // 5. WebSocketイベント処理のコールバックをメモ化
  // battlePhase.transitionToをRefで保持
  const battlePhaseTransitionRef = useRef(battlePhase.transitionTo);
  useEffect(() => {
    battlePhaseTransitionRef.current = battlePhase.transitionTo;
  }, [battlePhase.transitionTo]);

  const handlePhaseTransition = useCallback((newPhase: string) => {
    console.log(`[BattlePage] Transitioning to phase: ${newPhase}`);
    if (newPhase === "clap_time") {
      battlePhaseTransitionRef.current("clap_time_1");
    } else {
      battlePhaseTransitionRef.current(newPhase as BattlePhase);
    }
  }, []);

  // refetch関数をRefで保持
  const refetchRoomRef = useRef(refetchRoom);
  const refetchBattleRef = useRef(refetchBattle);
  const refetchImagesRef = useRef(refetchImages);

  useEffect(() => {
    refetchRoomRef.current = refetchRoom;
    refetchBattleRef.current = refetchBattle;
    refetchImagesRef.current = refetchImages;
  }, [refetchRoom, refetchBattle, refetchImages]);

  const handlePlayerChange = useCallback(() => {
    console.log("[BattlePage] Player change detected");
    refetchRoomRef.current();
    refetchBattleRef.current();
  }, []);

  const handleImageUpdate = useCallback(async () => {
    console.log("[BattlePage] Image update detected");
    await Promise.all([refetchBattleRef.current(), refetchImagesRef.current()]);
    console.log("[BattlePage] Images and battle info refetched");
  }, []);

  const handleClapUpdate = useCallback(() => {
    console.log("[BattlePage] Clap update detected");
    const id = `remote-clap-${Date.now()}-${Math.random()}`;
    const offsetX = Math.random() * 60 - 30; // -30px ~ +30px
    setRemoteClapEffects((prev) => [
      ...prev,
      { id, timestamp: Date.now(), offsetX },
    ]);
    setTimeout(() => {
      setRemoteClapEffects((prev) => prev.filter((e) => e.id !== id));
    }, 1500);
  }, []);

  const handleResultStart = useCallback(async () => {
    console.log("[BattlePage] Result phase started, fetching battle result...");

    // AI採点が完了していない可能性があるため、リトライロジックを実装
    const maxRetries = 7;
    const retryDelay = 3000; // 3秒
    let lastResult: GetBattleResultResponse | null = null;

    for (let attempt = 1; attempt <= maxRetries; attempt++) {
      console.log(
        `[BattlePage] Attempt ${attempt}/${maxRetries} to fetch result`
      );

      const result = await getResult(battleID);
      if (result) {
        lastResult = result;

        // AI採点が完了しているかチェック（ai_explanationが空でないこと）
        const isAIJudgingComplete = result.results.every(
          (r) => r.ai_explanation && r.ai_explanation.trim() !== ""
        );

        if (isAIJudgingComplete) {
          setBattleResult(result);
          console.log(
            "[BattlePage] Battle result fetched successfully:",
            result
          );
          return;
        } else {
          console.log(
            "[BattlePage] AI judging not complete yet, will retry..."
          );
        }
      }

      if (attempt < maxRetries) {
        console.log(
          `[BattlePage] Result not ready, waiting ${retryDelay}ms before retry...`
        );
        await new Promise((resolve) => setTimeout(resolve, retryDelay));
      }
    }

    // 最後のリトライでも完了しなかった場合、取得できた結果があればそれを使用
    if (lastResult) {
      setBattleResult(lastResult);
      console.warn(
        "[BattlePage] AI judging not fully complete, but proceeding with available result:",
        lastResult
      );
    } else {
      console.error(
        "[BattlePage] Failed to fetch battle result after all retries"
      );
    }
  }, [battleID, getResult]);

  // 6. WebSocketイベント処理（フェーズ遷移をトリガー）
  useBattleWebSocket({
    battleId: battleID,
    roomId: battle?.roomId || null,
    onPhaseTransition: handlePhaseTransition,
    onPlayerChange: handlePlayerChange,
    onImageUpdate: handleImageUpdate,
    onClapUpdate: handleClapUpdate,
    onResultStart: handleResultStart,
  });

  // 7. バトル情報取得後、途中参加を考慮してselectingフェーズに自動遷移
  // battle存在 = ゲーム開始済みと判定
  // 将来的に「全員揃うまで待機」が必要な場合は、バックエンドにバトル状態を追加して対応
  useEffect(() => {
    if (battle && battlePhase.phase === "waiting") {
      console.log(
        "[BattlePage] Battle loaded (途中参加), transitioning to selecting phase"
      );
      battlePhase.transitionTo("selecting");
      timer.resetTimer(PHASE_CONFIGS.selecting.duration);
      timer.startTimer();
    }
  }, [battle, timer, battlePhase]);

  // 8. 画像選択ロジック（フェーズに依存）
  const handleSendSuccess = useCallback(() => {
    console.log("[BattlePage] Image sent successfully");
  }, []);

  const handleSendError = useCallback(() => {
    console.log("[BattlePage] Failed to send image");
    setShowErrorDialog(true);
  }, []);

  const imageSelection = useImageSelection({
    battleId: battleID,
    userId: userId,
    canSelect: battlePhase.canSelectImage,
    onSendSuccess: handleSendSuccess,
    onSendError: handleSendError,
  });

  // 9. 拍手機能
  const { sendClap, isConnected: isClapConnected } = useWebSocketClap();

  // 現在の拍手ターゲットユーザーを計算
  const currentClapTarget = useMemo(() => {
    const clapIndex = CLAP_PHASES.indexOf(battlePhase.phase);
    if (clapIndex === -1) return null; // 拍手フェーズでない

    const targetPlayer = players[clapIndex];
    return targetPlayer || null;
  }, [battlePhase.phase, players]);

  // 拍手ハンドラー
  const handleClap = useCallback(() => {
    if (!currentClapTarget || !isClapConnected) return;

    try {
      sendClap({
        userId: userId,
        targetUserId: currentClapTarget.id,
        battleId: battleID,
        count: 1,
      });
      console.log(`[BattlePage] Clap sent to ${currentClapTarget.name}`);
    } catch (error) {
      console.error("[BattlePage] Failed to send clap:", error);
    }
  }, [currentClapTarget, isClapConnected, sendClap, userId, battleID]);

  return (
    <Battle
      // フェーズ情報
      phase={battlePhase.phase}
      phaseMessage={toastMessage}
      showPhaseMessage={showToast}
      // タイマー関連
      timeLeft={timer.timeLeft}
      isWarning={timer.isWarning}
      // 画像選択関連
      selectedImage={imageSelection.selectedImage}
      displayedImage={displayedImage}
      isDragging={imageSelection.isDragging}
      isImageSent={imageSelection.isImageSent}
      isSending={imageSelection.isSending}
      isCompressing={imageSelection.isCompressing}
      canSelect={battlePhase.canSelectImage}
      fileInputRef={imageSelection.fileInputRef}
      onImageSelect={imageSelection.handleImageSelect}
      onOpenAlbum={imageSelection.handleOpenAlbum}
      onCancel={imageSelection.handleCancel}
      onConfirmImage={imageSelection.handleConfirmImage}
      onDragEnter={imageSelection.handleDragEnter}
      onDragLeave={imageSelection.handleDragLeave}
      onDragOver={imageSelection.handleDragOver}
      onDrop={imageSelection.handleDrop}
      // バトル情報
      theme={battle?.theme || null}
      isLoading={battleInfoLoading}
      error={battleInfoError}
      players={players}
      images={images}
      remoteClapEffects={remoteClapEffects}
      // 拍手機能
      canClap={battlePhase.canClap}
      onClap={handleClap}
      // 結果情報
      battleResult={battleResult}
      battleId={battleID}
      // エラーダイアログ
      showErrorDialog={showErrorDialog}
      onCloseErrorDialog={() => setShowErrorDialog(false)}
    />
  );
}
