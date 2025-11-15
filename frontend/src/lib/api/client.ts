/**
 * API Client - 共通のHTTPリクエスト処理
 */

// fetchApi関数のオプション型
// 標準のfetchオプション + クエリパラメータ用のparams
type FetchOptions = RequestInit & {
  params?: Record<string, string | number>;
};

/**
 * API呼び出しエラー用のクラス
 * HTTPステータスコードとデバッグ情報を保持
 */
class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
    public data?: unknown
  ) {
    super(message);
    this.name = "ApiError";
  }
}

/**
 * API リクエストを実行する共通関数
 *
 * @template T - レスポンスの型
 * @param endpoint - APIのエンドポイント（例: '/room/info'）
 * @param options - リクエストオプション
 * @returns パースされたレスポンス
 * @throws {ApiError} HTTPエラーの場合
 */
export async function fetchApi<T>(
  endpoint: string,
  options: FetchOptions = {}
): Promise<T> {
  // paramsを分離（クエリパラメータ用）
  const { params, ...init } = options;

  // URLを構築
  let url = `${process.env.NEXT_PUBLIC_API_URL}${endpoint}`;

  // クエリパラメータがあればURLに追加
  if (params) {
    const queryString = new URLSearchParams(
      Object.entries(params).map(([key, value]) => [key, String(value)])
    ).toString();
    url += `?${queryString}`;
  }

  // HTTPリクエストを送信
  const response = await fetch(url, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...init.headers,
    },
  });

  // エラーレスポンスの場合は例外をスロー
  if (!response.ok) {
    const errorText = await response.text();
    throw new ApiError(response.status, errorText, { url, init });
  }

  // 204 No Content の場合は空オブジェクトを返す
  if (response.status === 204) {
    return {} as T;
  }

  // JSONをパースして返す
  return response.json();
}
