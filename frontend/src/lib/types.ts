/**
 * デバッグ用の型定義
 */
export type DebugMessage = {
  type?: string;
  api?: string;
  payload?: unknown;
  response?: unknown;
  timestamp: string;
} | null;
