import { useEffect } from "react";
import { useWebSocketEvents } from "@/src/lib/websocket/hooks/useWebSocketEvents";
import type {
  ImageSendPayload,
  PlayerJoinRoomPayload,
  PlayerLeaveRoomPayload,
  StartGamePayload,
  StartClapTimePayload,
  ClapSendPayload,
} from "@/src/lib/websocket/types";
import type { BattlePhase } from "./useBattlePhase";

type UseBattleWebSocketOptions = {
  battleId: string;
  roomId: string | null;
  onPhaseTransition: (phase: BattlePhase) => void;
  onPlayerChange: () => void;
  onImageUpdate: () => void;
  onClapUpdate?: () => void;
};

/**
 * バトルWebSocketイベント処理フック
 * WebSocketイベントの購読とフェーズ遷移のトリガーを担当
 */
export function useBattleWebSocket(options: UseBattleWebSocketOptions) {
  const { subscribe } = useWebSocketEvents();
  const {
    battleId,
    roomId,
    onPhaseTransition,
    onPlayerChange,
    onImageUpdate,
    onClapUpdate,
  } = options;

  useEffect(() => {
    if (!roomId) return;

    console.log(
      `[useBattleWebSocket] Subscribing to WebSocket events for room: ${roomId}`
    );

    const unsubscribers = [
      // ゲーム開始 → 選択フェーズ
      subscribe("start_game", (payload: StartGamePayload) => {
        console.log("[useBattleWebSocket] Game started:", payload);
        onPhaseTransition("selecting");
      }),

      // 拍手タイム開始 → 拍手フェーズ
      subscribe("start_clap_time", (payload: StartClapTimePayload) => {
        console.log("[useBattleWebSocket] Clap time started:", payload);
        onPhaseTransition("clap_time");
      }),

      // プレイヤー参加
      subscribe("player_join_room", (payload: PlayerJoinRoomPayload) => {
        console.log("[useBattleWebSocket] Player joined:", payload);
        onPlayerChange();
      }),

      // プレイヤー退出
      subscribe("player_leave_room", (payload: PlayerLeaveRoomPayload) => {
        console.log("[useBattleWebSocket] Player left:", payload);
        onPlayerChange();
      }),

      // 画像送信
      subscribe("image_send", (payload: ImageSendPayload) => {
        console.log("[useBattleWebSocket] Image sent:", payload);
        onImageUpdate();
      }),

      // 拍手送信
      subscribe("clap_send", (payload: ClapSendPayload) => {
        console.log("[useBattleWebSocket] Clap sent:", payload);
        onClapUpdate?.();
      }),
    ];

    return () => {
      console.log(
        `[useBattleWebSocket] Unsubscribing from WebSocket events for room: ${roomId}`
      );
      unsubscribers.forEach((unsub) => unsub());
    };
  }, [
    roomId,
    subscribe,
    onPhaseTransition,
    onPlayerChange,
    onImageUpdate,
    onClapUpdate,
    battleId,
  ]);
}
