"use client";

import { useState, useEffect, useCallback } from "react";
import { useBattleInfo } from "@/src/hooks/useBattleInfo";
import { useRoomInfo } from "@/src/hooks/useRoomInfo";
import { useBattlePhase, PHASE_CONFIGS } from "@/src/hooks/useBattlePhase";
import { useBattleTimer } from "@/src/hooks/useBattleTimer";
import { useBattleWebSocket } from "@/src/hooks/useBattleWebSocket";
import { useImageSelection } from "@/src/hooks/useImageSelection";
import { Battle } from "../components/Battle";

type BattlePageProps = {
  battleID: string;
};

// TODO: 実際のユーザーIDを取得する仕組みが必要
const MOCK_USER_ID = "550e8400-e29b-41d4-a716-446655440001";

/**
 * バトル画面のコンテナコンポーネント (Container)
 * ロジック・状態管理を担当
 */
export function BattlePage({ battleID }: BattlePageProps) {
  const [showErrorDialog, setShowErrorDialog] = useState(false);

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

  // 3. フェーズ管理（WebSocket駆動）
  const battlePhase = useBattlePhase("waiting");

  // 4. タイマー（表示のみ）
  const timer = useBattleTimer({
    initialTime: battlePhase.config.duration,
    autoStart: false,
    onWarning: useCallback((secondsLeft: number) => {
      console.log(`[BattlePage] Warning: ${secondsLeft} seconds left`);
    }, []),
    onTimeUp: useCallback(() => {
      console.log("[BattlePage] Time is up!");
    }, []),
  });

  // 5. WebSocketイベント処理のコールバックをメモ化
  const handlePhaseTransition = useCallback(
    (
      newPhase: "waiting" | "selecting" | "clap_time" | "result" | "finished"
    ) => {
      console.log(`[BattlePage] Transitioning to phase: ${newPhase}`);
      battlePhase.transitionTo(newPhase);
      // フェーズ遷移時にタイマーをリセット
      const config = PHASE_CONFIGS[newPhase];
      timer.resetTimer(config.duration);
      timer.startTimer();
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [battlePhase.transitionTo, timer.resetTimer, timer.startTimer]
  );

  const handlePlayerChange = useCallback(() => {
    console.log("[BattlePage] Player change detected");
    refetchRoom();
    refetchBattle();
  }, [refetchRoom, refetchBattle]);

  const handleImageUpdate = useCallback(() => {
    console.log("[BattlePage] Image update detected");
    refetchBattle();
    refetchImages();
  }, [refetchBattle, refetchImages]);

  const handleClapUpdate = useCallback(() => {
    console.log("[BattlePage] Clap update detected");
    // TODO: 拍手を受け取ったら演出や音を鳴らす？
    // スコアのフェッチは最後で良さげ
  }, []);

  // 6. WebSocketイベント処理（フェーズ遷移をトリガー）
  useBattleWebSocket({
    battleId: battleID,
    roomId: battle?.roomId || null,
    onPhaseTransition: handlePhaseTransition,
    onPlayerChange: handlePlayerChange,
    onImageUpdate: handleImageUpdate,
    onClapUpdate: handleClapUpdate,
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
  }, [battle, battlePhase, timer]);

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
    userId: MOCK_USER_ID,
    canSelect: battlePhase.canSelectImage,
    onSendSuccess: handleSendSuccess,
    onSendError: handleSendError,
  });

  return (
    <Battle
      // フェーズ情報
      phase={battlePhase.phase}
      // タイマー関連
      timeLeft={timer.timeLeft}
      isWarning={timer.isWarning}
      // 画像選択関連
      selectedImage={imageSelection.selectedImage}
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
      players={room?.users}
      images={images}
      // エラーダイアログ
      showErrorDialog={showErrorDialog}
      onCloseErrorDialog={() => setShowErrorDialog(false)}
    />
  );
}
