"use client";

import { createContext, useContext, useState } from "react";

type Player = {
  id: string;
  name: string;
};

type Room = {
  roomId: string;
  players: Player[];
};

type RoomContextType = {
  rooms: Room[];
  joinRoom: (roomId: string, playerName: string) => void;
};

const RoomContext = createContext<RoomContextType | null>(null);

export function RoomProvider({ children }: { children: React.ReactNode }) {
  const [rooms, setRooms] = useState<Room[]>([]);

  const joinRoom = (roomId: string, playerName: string) => {
    setRooms((prev) => {
      const existed = prev.find((r) => r.roomId === roomId);

      // 部屋が既にある → プレイヤー追加
      if (existed) {
        return prev.map((r) =>
          r.roomId === roomId
            ? {
                ...r,
                players: [
                  ...r.players,
                  { id: crypto.randomUUID(), name: playerName },
                ],
              }
            : r
        );
      }

      // なければ新規作成
      return [
        ...prev,
        {
          roomId,
          players: [{ id: crypto.randomUUID(), name: playerName }],
        },
      ];
    });
  };

  return (
    <RoomContext.Provider value={{ rooms, joinRoom }}>
      {children}
    </RoomContext.Provider>
  );
}

export const useRoomStore = () => useContext(RoomContext)!;
