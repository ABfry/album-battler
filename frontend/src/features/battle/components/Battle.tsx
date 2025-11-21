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
  // エラーダイアログ
  showErrorDialog,
  onCloseErrorDialog,
}: BattleProps) {
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
