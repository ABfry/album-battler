import Cookies from "js-cookie";

/**
 * クライアントサイドでCookieからユーザーIDを取得する
 * @returns ユーザーID（存在しない場合はnull）
 */
export function getUserIdClient(): string | null {
  return Cookies.get("userId") || null;
}
