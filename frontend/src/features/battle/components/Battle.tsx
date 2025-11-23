import Image from "next/image";
import { useState, useCallback, useEffect } from "react";
import type { CSSProperties } from "react";
import { useRouter } from "next/navigation";
import { motion, AnimatePresence } from "framer-motion";
import type {
  UserInfo,
  BattleImage,
  GetBattleResultResponse,
} from "@/src/lib/api/types";
import type { BattlePhase } from "@/src/hooks/useBattlePhase";
import { ImageFrame } from "./ImageFrame";
import { PlayerList } from "./PlayerList";
import { ResultDisplay } from "./ResultDisplay";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/src/components/ui/alert-dialog";

export type ClapEffect = {
  id: string;
  timestamp: number;
  offsetX: number; // 演出用の横移動量(px)
  originXPercent?: number; // 発生位置（親幅に対する%）
};

type BattleProps = {
  // フェーズ情報
  phase: BattlePhase;
  phaseMessage: string;
  showPhaseMessage: boolean;
  // タイマー関連
  timeLeft: number;
  isWarning: boolean;

  // 画像選択関連
  selectedImage: string | null;
  displayedImage?: string | null; // 拍手フェーズで表示する画像
  isDragging: boolean;
  isImageSent: boolean;
  isSending: boolean;
  isCompressing: boolean;
  canSelect: boolean; // 画像選択可能か
  fileInputRef: React.RefObject<HTMLInputElement | null>;
  onImageSelect: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onOpenAlbum: () => void;
  onCancel: () => void;
  onConfirmImage: () => void;
  onDragEnter: (e: React.DragEvent) => void;
  onDragLeave: (e: React.DragEvent) => void;
  onDragOver: (e: React.DragEvent) => void;
  onDrop: (e: React.DragEvent) => void;

  // バトル情報
  theme: string | null;
  isLoading: boolean;
  error: string | null;
  showThemeIntro: boolean;
  players: UserInfo[] | undefined;
  images: BattleImage[];
  remoteClapEffects: ClapEffect[];

  // 拍手機能
  canClap: boolean;
  onClap: () => void;

  // 結果情報
  battleResult: GetBattleResultResponse | null;
  battleId: string;

  // エラーダイアログ
  showErrorDialog: boolean;
  onCloseErrorDialog: () => void;
};

/**
 * バトル画面のメインコンポーネント (Presentational)
 */
export function Battle({
  // phase は将来的に使用予定のため型定義のみ保持
  phase,
  // タイマー関連
  timeLeft,
  isWarning,
  // 画像選択関連
  selectedImage,
  displayedImage,
  isDragging,
  isImageSent,
  isSending,
  isCompressing,
  canSelect,
  fileInputRef,
  onImageSelect,
  onOpenAlbum,
  onCancel,
  onConfirmImage,
  onDragEnter,
  onDragLeave,
  onDragOver,
  onDrop,
  // バトル情報
  theme,
  isLoading,
  error,
  phaseMessage,
  showPhaseMessage,
  players,
  showThemeIntro,
  images,
  remoteClapEffects,
  // 拍手機能
  canClap,
  onClap,
  // 結果情報
  battleResult,
  battleId,
  // エラーダイアログ
  showErrorDialog,
  onCloseErrorDialog,
}: BattleProps) {
  // 拍手エフェクトの管理
  const [clapEffects, setClapEffects] = useState<ClapEffect[]>([]);

  // 全てのエフェクトに対するタイマーをセットアップ＆クリーンアップ
  useEffect(() => {
    const timers = clapEffects.map((effect) => {
      return setTimeout(() => {
        setClapEffects((prev) => prev.filter((e) => e.id !== effect.id));
      }, 2000);
    });

    // クリーンアップ: コンポーネントのアンマウント時または clapEffects 変更時にタイマーをクリア
    return () => {
      timers.forEach(clearTimeout);
    };
  }, [clapEffects]);

  // 拍手ボタンクリック時のハンドラー
  const handleClapClick = useCallback(() => {
    // 元の拍手処理を実行
    onClap();

    // 拍手エフェクトを追加（ランダム値はここで計算）
    const newEffect: ClapEffect = {
      id: `clap-${Date.now()}-${Math.random()}`,
      timestamp: Date.now(),
      offsetX: Math.random() * 60 - 30, // -30px ~ +30px のランダムな横移動
      originXPercent: 10 + Math.random() * 80, // ボタン幅ほぼ全体をカバーする発生位置
    };
    setClapEffects((prev) => [...prev, newEffect]);
  }, [onClap]);

  const router = useRouter();

  const renderQuestionHeader = (variant: "large" | "small") => {
    const isLarge = variant === "large";
    const maxWidthClass = isLarge
      ? "w-full px-4 sm:px-6 md:px-0 max-w-[1200px]"
      : "w-full px-4 sm:px-6 max-w-[1500px]";
    const frameStyle: CSSProperties = {
      aspectRatio: "260 / 72",
      minHeight: isLarge ? 140 : 160,
    };
    const textStyle: CSSProperties = {
      top: "33%",
      bottom: 0,
      padding: isLarge ? "0 10%" : "0 10%",
      maxHeight: "67%",
      fontSize: isLarge ? "clamp(16px, 3vw, 24px)" : "clamp(18px, 3vw, 28px)",
      lineHeight: 1.2,
      wordBreak: "break-word",
      whiteSpace: "pre-wrap",
      textAlign: "center",
      display: "flex",
      alignItems: "center",
      justifyContent: "center",
    };

    return (
      <div
        className={`relative mx-auto w-full overflow-hidden ${maxWidthClass}`}
        style={frameStyle}
      >
        <Image
          src="/question-header.png"
          alt="Question Header"
          fill
          className="object-contain"
          priority
          sizes="300px"
        />
        <div
          className="absolute right-0 left-0 flex items-center justify-center px-3 text-center leading-tight font-black text-slate-900"
          style={textStyle}
        >
          {theme || "テーマ未設定"}
        </div>
      </div>
    );
  };

  // 結果フェーズの表示 + 10秒後に詳細ページへ遷移
  useEffect(() => {
    if (phase === "result" && battleResult) {
      const timer = setTimeout(() => {
        router.push(`/result/${battleId}`);
      }, 10000);

      return () => clearTimeout(timer);
    }
  }, [phase, battleResult, battleId, router]);

  // 結果フェーズの表示
  if (phase === "result" && battleResult) {
    return (
      <ResultDisplay battleResult={battleResult} players={players || []} />
    );
  }

  return (
    <div className="flex h-dvh flex-col items-center justify-center overflow-hidden p-4">
      <div className="relative flex h-full w-full max-w-4xl flex-col items-center gap-2 pt-16 pb-2">
        {showThemeIntro && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70">
            <div className="flex flex-col items-center gap-6 px-8 py-6">
              {renderQuestionHeader("large")}
            </div>
          </div>
        )}
        {/* タイマー表示（右上固定） */}
        <div className="absolute top-4 right-4 z-30 md:top-6 md:right-6">
          <div
            className={`inline-block rounded-lg px-6 py-1.5 text-3xl font-bold shadow-lg ${
              isWarning ? "bg-red-500 text-white" : "bg-white text-slate-800"
            }`}
          >
            {timeLeft}秒
          </div>
        </div>

        {/* お題表示（タイマーの下） */}
        <div className="flex w-full shrink-0 items-center justify-center">
          {isLoading ? (
            <span className="text-2xl font-black text-gray-400 md:text-3xl">
              Loading...
            </span>
          ) : error ? (
            <span className="text-2xl font-black text-red-500 md:text-3xl">
              Error: {error}
            </span>
          ) : showThemeIntro ? null : (
            <motion.div
              key={
                showThemeIntro ? "question-header-hidden" : "question-header"
              }
              initial={{ opacity: 0, scale: 0.9 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ duration: 0.3 }}
              className="mx-auto mb-8 flex items-center justify-center"
            >
              {renderQuestionHeader("small")}
            </motion.div>
          )}
        </div>

        {phaseMessage && showPhaseMessage && (
          <div className="pointer-events-none absolute inset-0 z-50 flex items-center justify-center">
            <div className="animate-in fade-in zoom-in rounded-3xl bg-black/70 px-8 py-6 text-4xl font-black text-white shadow-[0_0_30px_rgba(0,0,0,0.5)] ring-4 ring-white/30 duration-300">
              {phaseMessage}
            </div>
          </div>
        )}

        {/* バトル画像（中央） - min-h-0でflexboxの縮小を有効化 */}
        <div className="min-h-0 w-full flex-1">
          <ImageFrame
            selectedImage={selectedImage}
            displayedImage={displayedImage}
            isDragging={isDragging}
            isImageSent={isImageSent}
            isSending={isSending}
            isCompressing={isCompressing}
            canSelect={canSelect}
            fileInputRef={fileInputRef}
            onImageSelect={onImageSelect}
            onOpenAlbum={onOpenAlbum}
            onCancel={onCancel}
            onConfirmImage={onConfirmImage}
            onDragEnter={onDragEnter}
            onDragLeave={onDragLeave}
            onDragOver={onDragOver}
            onDrop={onDrop}
          />
        </div>

        {/* プレイヤー情報（最下部） */}
        <div className="w-full shrink-0">
          <PlayerList
            players={players}
            images={images}
            remoteClapEffects={remoteClapEffects}
          />
        </div>

        {/* 拍手ボタン（プレイヤーリスト下） */}
        {canClap && (
          <div className="relative mt-4 w-full">
            {/* 拍手エフェクト */}
            <div className="relative mx-auto w-full max-w-4xl">
              <AnimatePresence>
                {clapEffects.map((effect) => (
                  <motion.div
                    key={effect.id}
                    initial={{ y: 0, opacity: 1, scale: 1 }}
                    animate={{
                      y: -300,
                      opacity: 0,
                      x: effect.offsetX,
                    }}
                    exit={{ opacity: 0 }}
                    transition={{
                      duration: 2.0,
                      ease: "easeOut",
                    }}
                    className="pointer-events-none absolute bottom-full text-6xl"
                    style={{
                      left: `${effect.originXPercent ?? 50}%`,
                      transform: "translateX(-50%)",
                    }}
                  >
                    👏
                  </motion.div>
                ))}
              </AnimatePresence>

              {/* 拍手ボタン */}
              <button
                onClick={handleClapClick}
                className="mx-auto flex h-16 w-full max-w-md items-center justify-center rounded-2xl bg-linear-to-r from-yellow-400 to-orange-500 text-3xl font-bold shadow-2xl transition-transform hover:scale-[1.02] active:scale-95"
              >
                👏
              </button>
            </div>
          </div>
        )}
      </div>

      {/* エラーダイアログ */}
      <AlertDialog open={showErrorDialog} onOpenChange={onCloseErrorDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>送信エラー</AlertDialogTitle>
            <AlertDialogDescription>
              画像の送信に失敗しました。もう一度お試しください。
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogAction onClick={onCloseErrorDialog}>
              OK
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
