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
import { Battle } from "../components/Battle";
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
    return getUserIdClient() || "550e8400-e29b-41d4-a716-446655440001";
  });

  const [showErrorDialog, setShowErrorDialog] = useState(false);
  const [displayedImage, setDisplayedImage] = useState<string | null>(null);
  const [battleResult, setBattleResult] =
    useState<GetBattleResultResponse | null>(null);

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

  const playersRef = useRef(players);
  const imagesRef = useRef(images);

  useEffect(() => {
    playersRef.current = players;
  }, [players]);

  useEffect(() => {
    imagesRef.current = images;
  }, [images]);

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
          timer.resetTimer(60);
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
        onPhaseStart: async () => {
          console.log("結果発表");
          // 結果データを取得
          console.log("[BattlePage] Fetching battle result...");
          const result = await getResult(battleID);
          if (result) {
            setBattleResult(result);
            console.log("[BattlePage] Battle result fetched:", result);
          } else {
            console.error("[BattlePage] Failed to fetch battle result");
          }
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
    // TODO: 拍手を受け取ったら演出や音を鳴らす？
    // スコアのフェッチは最後で良さげ
  }, []);

  const handleResultStart = useCallback(async () => {
    console.log("[BattlePage] Fetching battle result...");
    const result = await getResult(battleID);
    if (result) {
      setBattleResult(result);
      console.log("[BattlePage] Battle result fetched:", result);
    } else {
      console.error("[BattlePage] Failed to fetch battle result");
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
      // 拍手機能
      canClap={battlePhase.canClap}
      onClap={handleClap}
      // 結果情報
      battleResult={battleResult}
      // エラーダイアログ
      showErrorDialog={showErrorDialog}
      onCloseErrorDialog={() => setShowErrorDialog(false)}
    />
  );
}
