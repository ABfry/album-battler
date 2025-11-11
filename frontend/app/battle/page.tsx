"use client";

import { useWebSocket } from "@/src/lib/websocket/contexts/WebSocketContext";
import Link from "next/link";

export default function BattlePage() {
  const { status, messages, sendMessage } = useWebSocket();

  const handleSendMessage = () => {
    sendMessage("Hello from Battle page!");
  };

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
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <h1 className="text-3xl font-bold text-black dark:text-white">
              Battle
            </h1>
            <div className="flex items-center gap-2">
              <div className={`h-3 w-3 rounded-full ${getStatusColor()}`}></div>
              <span className="text-sm text-zinc-600 dark:text-zinc-400">
                {status}
              </span>
            </div>
          </div>
          <Link
            href="/room"
            className="rounded-lg bg-zinc-800 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-zinc-700 dark:bg-zinc-700 dark:hover:bg-zinc-600"
          >
            Back to Room
          </Link>
        </div>

        <div className="rounded-lg border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-900">
          <h2 className="mb-4 text-xl font-semibold text-black dark:text-white">
            WebSocket Connection
          </h2>
          <p className="mb-4 text-zinc-600 dark:text-zinc-400">
            WebSocket接続は全ページで共有されています。
          </p>
          <button
            onClick={handleSendMessage}
            disabled={status !== "connected"}
            className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-zinc-400 disabled:hover:bg-zinc-400"
          >
            Send Test Message
          </button>
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

        <div className="rounded-lg border border-blue-200 bg-blue-50 p-4 dark:border-blue-900 dark:bg-blue-950">
          <p className="text-sm text-blue-800 dark:text-blue-200">
            💡 このページと /room
            ページを行き来しても、WebSocket接続は維持されます！
          </p>
        </div>
      </main>
    </div>
  );
}
