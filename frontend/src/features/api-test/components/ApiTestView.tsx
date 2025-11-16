import { RoomCreateForm } from "./RoomCreateForm";
import { RoomJoinForm } from "./RoomJoinForm";
import { RoomInfo } from "./RoomInfo";
import { RoomLeaveButton } from "./RoomLeaveButton";
import { GameStartButton } from "./GameStartButton";
import { BattleIDInfo } from "./BattleIDInfo";
import { BattleCreateForm } from "./BattleCreateForm";
import { BattleInfo } from "./BattleInfo";
import { ImageSendForm } from "./ImageSendForm";
import { ImageList } from "./ImageList";
import type {
  CreateRoomResponse,
  CreateBattleResponse,
  Room,
  Battle,
  BattleImage,
} from "@/src/lib/api/types";

type ApiTestViewProps = {
  createdRoom: CreateRoomResponse | null;
  room: Room | null;
  loading: boolean;
  roomLoading: boolean;
  error: string | null;
  roomError: string | null;
  joinSuccess: boolean;
  leaveSuccess: boolean;
  startSuccess: boolean;
  battleID: string | null;
  createdBattle: CreateBattleResponse | null;
  battle: Battle | null;
  images: BattleImage[];
  battleLoading: boolean;
  battleInfoLoading: boolean;
  battleError: string | null;
  battleInfoError: string | null;
  sendImageSuccess: boolean;
  onCreateRoom: (userId: string) => Promise<void>;
  onJoinRoom: (userId: string, roomNumber: number) => Promise<void>;
  onLeaveRoom: (roomId: string, userId: string) => Promise<void>;
  onStartGame: () => Promise<void>;
  onGetBattleID: () => Promise<void>;
  onCreateBattle: () => Promise<void>;
  onSendImage: (userId: string, imageBase64: string) => Promise<void>;
  onRefetch: () => void;
  onRefetchBattle: () => void;
  onRefetchImages: () => void;
};

/**
 * API テストページ Presentational Component
 * UIの表示のみを担当
 */
export function ApiTestView({
  createdRoom,
  room,
  loading,
  roomLoading,
  error,
  roomError,
  joinSuccess,
  leaveSuccess,
  startSuccess,
  battleID,
  createdBattle,
  battle,
  images,
  battleLoading,
  battleInfoLoading,
  battleError,
  battleInfoError,
  sendImageSuccess,
  onCreateRoom,
  onJoinRoom,
  onLeaveRoom,
  onStartGame,
  onGetBattleID,
  onCreateBattle,
  onSendImage,
  onRefetch,
  onRefetchBattle,
  onRefetchImages,
}: ApiTestViewProps) {
  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="mx-auto max-w-6xl">
        <h1 className="mb-8 text-3xl font-bold">REST API Test Page</h1>

        <div className="mb-6 rounded-lg bg-blue-50 p-4">
          <h2 className="mb-2 font-semibold">📝 使い方</h2>
          <ol className="list-inside list-decimal space-y-1 text-sm">
            <li>「Create Room」で部屋を作成</li>
            <li>「Join Room」で別のユーザーが参加（Room Numberを確認）</li>
            <li>「Room Info」で参加者を確認</li>
            <li>「Leave Room」で部屋から退出</li>
            <li>「Start Game」でゲームを開始（ホストのみ実行可能）</li>
            <li>「Get Battle ID」でバトルIDを取得（Start Game後に自動実行）</li>
            <li>「Create Battle」でバトルを作成（デバッグ用）</li>
            <li>「Battle Info」でバトル情報を確認</li>
            <li>「Send Image」で画像を送信</li>
            <li>「Battle Images」で送信された画像一覧を確認</li>
          </ol>
        </div>

        <div className="grid gap-6 md:grid-cols-2">
          {/* 1. 部屋作成 */}
          <RoomCreateForm
            onCreateRoom={onCreateRoom}
            loading={loading}
            error={error}
            result={createdRoom}
          />

          {/* 2. 部屋参加 */}
          <RoomJoinForm
            onJoinRoom={onJoinRoom}
            loading={loading}
            error={error}
            success={joinSuccess}
            autoFillRoomNumber={createdRoom?.room_number}
          />

          {/* 3. 部屋情報 */}
          <RoomInfo
            room={room}
            loading={roomLoading}
            error={roomError}
            onRefetch={onRefetch}
          />

          {/* 4. 部屋退出 */}
          <RoomLeaveButton
            roomId={createdRoom?.room_id || null}
            userId="550e8400-e29b-41d4-a716-446655440002"
            onLeaveRoom={onLeaveRoom}
            loading={loading}
            error={error}
            success={leaveSuccess}
          />

          {/* 5. ゲーム開始 */}
          <GameStartButton
            roomId={createdRoom?.room_id || null}
            userId="550e8400-e29b-41d4-a716-446655440001"
            onStartGame={onStartGame}
            loading={loading}
            error={error}
            success={startSuccess}
          />

          {/* 6. バトルID取得 */}
          <BattleIDInfo
            roomId={createdRoom?.room_id || null}
            battleID={battleID}
            onGetBattleID={onGetBattleID}
            loading={loading}
            error={error}
          />

          {/* 7. バトル作成 */}
          <BattleCreateForm
            roomId={createdRoom?.room_id || null}
            onCreateBattle={onCreateBattle}
            loading={battleLoading}
            error={battleError}
            result={createdBattle}
          />

          {/* 8. バトル情報 */}
          <BattleInfo
            battle={battle}
            loading={battleInfoLoading}
            error={battleInfoError}
            onRefetch={onRefetchBattle}
          />

          {/* 9. 画像送信 */}
          <ImageSendForm
            battleId={createdBattle?.BattleID || null}
            onSendImage={onSendImage}
            loading={battleLoading}
            error={battleError}
            success={sendImageSuccess}
          />

          {/* 10. 画像一覧 */}
          <ImageList
            images={images}
            loading={battleInfoLoading}
            error={battleInfoError}
            onRefetch={onRefetchImages}
          />
        </div>

        {/* デバッグ情報 */}
        <div className="mt-8 rounded-lg border bg-white p-4">
          <h2 className="mb-2 text-lg font-bold">🔍 Debug Info</h2>
          <pre className="overflow-x-auto rounded bg-gray-100 p-3 text-xs">
            {JSON.stringify(
              {
                createdRoom,
                room,
                createdBattle,
                battle,
                images,
                loading,
                roomLoading,
                battleLoading,
                battleInfoLoading,
                error,
                roomError,
                battleError,
                battleInfoError,
              },
              null,
              2
            )}
          </pre>
        </div>
      </div>
    </div>
  );
}
