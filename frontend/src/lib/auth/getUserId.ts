import { cookies } from "next/headers";

/**
 * サーバーサイドでCookieからユーザーIDを取得する
 * @returns ユーザーID（存在しない場合はnull）
 */
export async function getUserId(): Promise<string | null> {
  const cookieStore = await cookies();
  const userIdCookie = cookieStore.get("userId");
  return userIdCookie?.value || null;
}
