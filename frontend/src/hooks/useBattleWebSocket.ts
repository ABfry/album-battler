import { useEffect } from "react";
import { useWebSocketEvents } from "@/src/lib/websocket/hooks/useWebSocketEvents";
import type {
  ImageSendPayload,
  PlayerJoinRoomPayload,
  PlayerLeaveRoomPayload,
  StartGamePayload,
  StartClapTimePayload,
} from "@/src/lib/websocket/types";
import type { BattlePhase } from "./useBattlePhase";

type UseBattleWebSocketOptions = {
  battleId: string;
  roomId: string | null;
  onPhaseTransition: (phase: BattlePhase) => void;
  onPlayerChange: () => void;
  onImageUpdate: () => void;
};

/**
 * バトルWebSocketイベント処理フック
 * WebSocketイベントの購読とフェーズ遷移のトリガーを担当
 */
export function useBattleWebSocket(options: UseBattleWebSocketOptions) {
  const { subscribe } = useWebSocketEvents();

  useEffect(() => {
    if (!options.roomId) return;

    console.log(
      `[useBattleWebSocket] Subscribing to WebSocket events for room: ${options.roomId}`
    );

    const unsubscribers = [
      // ゲーム開始 → 選択フェーズ
      subscribe("start_game", (payload: StartGamePayload) => {
        console.log("[useBattleWebSocket] Game started:", payload);
        options.onPhaseTransition("selecting");
      }),

      // 拍手タイム開始 → 拍手フェーズ
      subscribe("start_clap_time", (payload: StartClapTimePayload) => {
        console.log("[useBattleWebSocket] Clap time started:", payload);
        options.onPhaseTransition("clap_time");
      }),

      // プレイヤー参加
      subscribe("player_join_room", (payload: PlayerJoinRoomPayload) => {
        console.log("[useBattleWebSocket] Player joined:", payload);
        options.onPlayerChange();
      }),

      // プレイヤー退出
      subscribe("player_leave_room", (payload: PlayerLeaveRoomPayload) => {
        console.log("[useBattleWebSocket] Player left:", payload);
        options.onPlayerChange();
      }),

      // 画像送信
      subscribe("image_send", (payload: ImageSendPayload) => {
        console.log("[useBattleWebSocket] Image sent:", payload);
        options.onImageUpdate();
      }),
    ];

    return () => {
      console.log(
        `[useBattleWebSocket] Unsubscribing from WebSocket events for room: ${options.roomId}`
      );
      unsubscribers.forEach((unsub) => unsub());
    };
  }, [options.roomId, subscribe, options]);
}
