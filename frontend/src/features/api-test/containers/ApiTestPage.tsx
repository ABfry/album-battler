"use client";

import { useState, useEffect, useRef } from "react";
import { useRoom } from "@/src/hooks/useRoom";
import { useRoomInfo } from "@/src/hooks/useRoomInfo";
import { useBattle } from "@/src/hooks/useBattle";
import { useBattleInfo } from "@/src/hooks/useBattleInfo";
import { useWebSocket } from "@/src/lib/websocket/contexts/WebSocketContext";
import { useWebSocketEvents } from "@/src/lib/websocket/hooks/useWebSocketEvents";
import { ApiTestView } from "../components/ApiTestView";
import type {
  CreateRoomResponse,
  CreateBattleResponse,
} from "@/src/lib/api/types";
import type {
  PlayerJoinRoomPayload,
  PlayerLeaveRoomPayload,
  StartGamePayload,
} from "@/src/lib/websocket/types";

/**
 * API テストページ Container
 * ロジック・状態管理を担当
 */
export function ApiTestPage() {
  // 部屋作成・参加・退出・ゲーム開始の操作
  const {
    createRoom,
    joinRoom,
    leaveRoom,
    startGame,
    getBattleID,
    loading,
    error,
  } = useRoom();

  // バトル作成・画像送信の操作
  const {
    createBattle,
    sendImage,
    loading: battleLoading,
    error: battleError,
  } = useBattle();

  // WebSocket接続
  const { connect, disconnect } = useWebSocket();
  const { subscribe } = useWebSocketEvents();

  // 作成された部屋の情報を保持
  const [createdRoom, setCreatedRoom] = useState<CreateRoomResponse | null>(
    null
  );

  // 作成されたバトルの情報を保持
  const [createdBattle, setCreatedBattle] =
    useState<CreateBattleResponse | null>(null);

  // 部屋情報の取得（作成後に自動取得）
  const {
    room,
    loading: roomLoading,
    error: roomError,
    refetch,
  } = useRoomInfo(createdRoom?.room_id || null);

  // バトル情報の取得
  const {
    battle,
    images,
    loading: battleInfoLoading,
    error: battleInfoError,
    refetch: refetchBattle,
    refetchImages,
  } = useBattleInfo(createdBattle?.BattleID || null);

  // 各操作の成功状態
  const [joinSuccess, setJoinSuccess] = useState(false);
  const [leaveSuccess, setLeaveSuccess] = useState(false);
  const [startSuccess, setStartSuccess] = useState(false);
  const [sendImageSuccess, setSendImageSuccess] = useState(false);
  const [battleID, setBattleID] = useState<string | null>(null);

  // デバッグ用: 最新のWebSocketメッセージとAPIレスポンス
  const [latestWsMessage, setLatestWsMessage] = useState<any>(null);
  const [latestApiResponse, setLatestApiResponse] = useState<any>(null);

  const createdRoomRef = useRef(createdRoom);
  const getBattleIDRef = useRef(getBattleID);

  useEffect(() => {
    createdRoomRef.current = createdRoom;
  }, [createdRoom]);

  useEffect(() => {
    getBattleIDRef.current = getBattleID;
  }, [getBattleID]);

  // WebSocket接続の初期化
  useEffect(() => {
    connect();
    return () => disconnect();
  }, [connect, disconnect]);

  // WebSocketイベントの購読
  useEffect(() => {
    const unsubscribeJoin = subscribe(
      "player_join_room",
      (payload: PlayerJoinRoomPayload) => {
        console.log("Player joined room:", payload.room_id);
        setLatestWsMessage({
          type: "player_join_room",
          payload,
          timestamp: new Date().toISOString(),
        });
        refetch();
      }
    );

    const unsubscribeLeave = subscribe(
      "player_leave_room",
      (payload: PlayerLeaveRoomPayload) => {
        console.log("Player left room:", payload.room_id);
        setLatestWsMessage({
          type: "player_leave_room",
          payload,
          timestamp: new Date().toISOString(),
        });
        refetch();
      }
    );

    const unsubscribeStart = subscribe(
      "start_game",
      async (payload: StartGamePayload) => {
        console.log("Game started:", payload.room_id);
        setLatestWsMessage({
          type: "start_game",
          payload,
          timestamp: new Date().toISOString(),
        });
        refetch();
        // ゲーム開始時に自動でBattle IDを取得
        if (createdRoomRef.current?.room_id === payload.room_id) {
          const id = await getBattleIDRef.current(payload.room_id);
          setBattleID(id);
        }
      }
    );

    return () => {
      unsubscribeJoin();
      unsubscribeLeave();
      unsubscribeStart();
    };
  }, [subscribe, refetch]);

  // 1. 部屋作成
  const handleCreateRoom = async (userId: string) => {
    const result = await createRoom(userId);
    if (result) {
      setCreatedRoom(result);
      setJoinSuccess(false);
      setStartSuccess(false);
      setLatestApiResponse({
        api: "createRoom",
        response: result,
        timestamp: new Date().toISOString(),
      });
    }
  };

  // 2. 部屋参加
  const handleJoinRoom = async (userId: string, roomNumber: number) => {
    const success = await joinRoom(userId, roomNumber);
    setJoinSuccess(success);
    setLatestApiResponse({
      api: "joinRoom",
      response: { success },
      timestamp: new Date().toISOString(),
    });
  };

  // 3. 部屋退出
  const handleLeaveRoom = async (roomId: string, userId: string) => {
    const success = await leaveRoom(roomId, userId);
    setLeaveSuccess(success);
    setLatestApiResponse({
      api: "leaveRoom",
      response: { success },
      timestamp: new Date().toISOString(),
    });
    if (success) {
      setTimeout(() => refetch(), 500);
    }
  };

  // 4. ゲーム開始
  const handleStartGame = async () => {
    if (!createdRoom) return;
    const success = await startGame(
      createdRoom.room_id,
      "550e8400-e29b-41d4-a716-446655440001"
    );
    setStartSuccess(success);
    setLatestApiResponse({
      api: "startGame",
      response: { success },
      timestamp: new Date().toISOString(),
    });
  };

  // 5. バトルID取得
  const handleGetBattleID = async () => {
    if (!createdRoom) return;
    const id = await getBattleID(createdRoom.room_id);
    setBattleID(id);
    setLatestApiResponse({
      api: "getBattleID",
      response: { battle_id: id },
      timestamp: new Date().toISOString(),
    });
  };

  // 6. バトル作成
  const handleCreateBattle = async () => {
    if (!createdRoom) return;
    const result = await createBattle(createdRoom.room_id);
    if (result) {
      setCreatedBattle(result);
      setLatestApiResponse({
        api: "createBattle",
        response: result,
        timestamp: new Date().toISOString(),
      });
    }
  };

  // 7. 画像送信
  const handleSendImage = async (userId: string, imageBase64: string) => {
    if (!createdBattle) return;
    const success = await sendImage(
      createdBattle.BattleID,
      userId,
      imageBase64
    );
    setSendImageSuccess(success);
    setLatestApiResponse({
      api: "sendImage",
      response: { success },
      timestamp: new Date().toISOString(),
    });
    if (success) {
      // 画像送信成功後、画像一覧を再取得
      setTimeout(() => refetchImages(), 500);
    }
  };

  return (
    <ApiTestView
      createdRoom={createdRoom}
      room={room}
      loading={loading}
      roomLoading={roomLoading}
      error={error}
      roomError={roomError}
      joinSuccess={joinSuccess}
      leaveSuccess={leaveSuccess}
      startSuccess={startSuccess}
      battleID={battleID}
      createdBattle={createdBattle}
      battle={battle}
      images={images}
      battleLoading={battleLoading}
      battleInfoLoading={battleInfoLoading}
      battleError={battleError}
      battleInfoError={battleInfoError}
      sendImageSuccess={sendImageSuccess}
      latestWsMessage={latestWsMessage}
      latestApiResponse={latestApiResponse}
      onCreateRoom={handleCreateRoom}
      onJoinRoom={handleJoinRoom}
      onLeaveRoom={handleLeaveRoom}
      onStartGame={handleStartGame}
      onGetBattleID={handleGetBattleID}
      onCreateBattle={handleCreateBattle}
      onSendImage={handleSendImage}
      onRefetch={refetch}
      onRefetchBattle={refetchBattle}
      onRefetchImages={refetchImages}
    />
  );
}
