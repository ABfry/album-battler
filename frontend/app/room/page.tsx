"use client";

import { useWebSocket } from "../providers/WebSocketProvider";
import Link from "next/link";

export default function RoomPage() {
  const { status, messages } = useWebSocket();

  const getStatusColor = () => {
    switch (status) {
      case "connected":
        return "bg-green-500";
      case "connecting":
        return "bg-yellow-500";
      case "disconnected":
        return "bg-gray-500";
      case "error":
        return "bg-red-500";
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 dark:bg-black">
      <main className="flex w-full max-w-2xl flex-col gap-6 p-8">
        <div className="flex items-center gap-4">
          <h1 className="text-3xl font-bold text-black dark:text-white">
            Room
          </h1>
          <div className="flex items-center gap-2">
            <div className={`h-3 w-3 rounded-full ${getStatusColor()}`}></div>
            <span className="text-sm text-zinc-600 dark:text-zinc-400">
              {status}
            </span>
          </div>
        </div>

        <div className="rounded-lg border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-900">
          <h2 className="mb-4 text-xl font-semibold text-black dark:text-white">
            WebSocket Status
          </h2>
          <p className="text-zinc-600 dark:text-zinc-400">
            Connection to:{" "}
            <code className="rounded bg-zinc-100 px-2 py-1 dark:bg-zinc-800">
              ws://localhost:8080/ws
            </code>
          </p>
        </div>

        <div className="rounded-lg border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-900">
          <h2 className="mb-4 text-xl font-semibold text-black dark:text-white">
            Messages
          </h2>
          <div className="flex max-h-64 flex-col gap-2 overflow-y-auto">
            {messages.length === 0 ? (
              <p className="text-sm text-zinc-500 dark:text-zinc-500">
                No messages yet...
              </p>
            ) : (
              messages.map((msg, idx) => (
                <div
                  key={idx}
                  className="rounded border border-zinc-200 bg-zinc-50 p-2 text-sm text-zinc-800 dark:border-zinc-700 dark:bg-zinc-800 dark:text-zinc-200"
                >
                  {msg}
                </div>
              ))
            )}
          </div>
        </div>

        <Link
          href="/battle"
          className="flex h-12 items-center justify-center rounded-lg bg-blue-600 text-base font-medium text-white transition-colors hover:bg-blue-700"
        >
          Go to Battle
        </Link>
      </main>
    </div>
  );
}
