import { useState, useCallback, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import type { UserInfo, BattleImage } from "@/src/lib/api/types";
import type { BattlePhase } from "@/src/hooks/useBattlePhase";
import { ImageFrame } from "./ImageFrame";
import { PlayerList } from "./PlayerList";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/src/components/ui/alert-dialog";

type ClapEffect = {
  id: string;
  timestamp: number;
  offsetX: number; // ランダムな横移動量
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
  displayedImage: string | null; // 拍手フェーズで表示する画像
  isDragging: boolean;
  isImageSent: boolean;
  isCompressing: boolean;
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
  players: UserInfo[] | undefined;
  images: BattleImage[];

  // 拍手機能
  canClap: boolean;
  onClap: () => void;

  // エラーダイアログ
  showErrorDialog: boolean;
  onCloseErrorDialog: () => void;
};

/**
 * バトル画面のメインコンポーネント (Presentational)
 */
export function Battle({
  // phase は将来的に使用予定のため型定義のみ保持
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  phase,
  // タイマー関連
  timeLeft,
  isWarning,
  // 画像選択関連
  selectedImage,
  displayedImage,
  isDragging,
  isImageSent,
  isCompressing,
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
  images,
  // 拍手機能
  canClap,
  onClap,
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
    };
    setClapEffects((prev) => [...prev, newEffect]);
  }, [onClap]);

  return (
    <div className="flex min-h-screen flex-col items-center justify-center p-4">
      <div className="relative flex h-screen w-full max-w-4xl flex-col items-center justify-between py-8">
        {/* タイマー表示（右上） */}
        <div className="absolute top-8 right-8">
          <div
            className={`rounded-lg px-6 py-3 text-4xl font-bold shadow-lg ${
              isWarning ? "bg-red-500 text-white" : "bg-white text-slate-800"
            }`}
          >
            {timeLeft}秒
          </div>
        </div>

        {/* お題表示（最上部） */}
        <div className="w-full text-center">
          <h1 className="text-4xl font-black text-slate-800 md:text-5xl">
            {isLoading ? (
              <span className="text-gray-400">Loading...</span>
            ) : error ? (
              <span className="text-red-500">Error: {error}</span>
            ) : theme ? (
              theme
            ) : (
              <span className="text-gray-400">テーマ未設定</span>
            )}
          </h1>
        </div>

        {phaseMessage && showPhaseMessage && (
          <div className="pointer-events-none absolute inset-0 z-50 flex items-center justify-center">
            <div className="animate-in fade-in zoom-in rounded-3xl bg-black/70 px-8 py-6 text-4xl font-black text-white shadow-[0_0_30px_rgba(0,0,0,0.5)] ring-4 ring-white/30 duration-300">
              {phaseMessage}
            </div>
          </div>
        )}

        {/* バトル画像（中央） */}
        <ImageFrame
          selectedImage={selectedImage}
          displayedImage={displayedImage}
          isDragging={isDragging}
          isImageSent={isImageSent}
          isCompressing={isCompressing}
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

        {/* プレイヤー情報（最下部） */}
        <PlayerList players={players} images={images} />

        {/* 拍手ボタン（右下） */}
        {canClap && (
          <div className="absolute right-8 bottom-8">
            {/* 拍手エフェクト */}
            <AnimatePresence>
              {clapEffects.map((effect) => (
                <motion.div
                  key={effect.id}
                  initial={{ y: 0, opacity: 1, scale: 1 }}
                  animate={{
                    y: -300, // 150 → 300: 距離を2倍に
                    opacity: 0,
                    x: effect.offsetX, // 事前に計算されたランダムな横移動量
                  }}
                  exit={{ opacity: 0 }}
                  transition={{
                    duration: 2.0, // 1.5 → 2.0: より長く、よりダイナミックに
                    ease: "easeOut",
                  }}
                  className="pointer-events-none absolute text-6xl"
                  style={{
                    left: "15%",
                    bottom: "100%",
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
              className="relative flex h-24 w-24 items-center justify-center rounded-full bg-linear-to-br from-yellow-400 to-orange-500 text-5xl shadow-2xl transition-transform hover:scale-110 active:scale-95"
            >
              👏
            </button>
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
