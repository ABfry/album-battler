import Cookies from "js-cookie";

/**
 * クライアントサイドでCookieからユーザーIDを取得する
 * @returns ユーザーID（存在しない場合はnull）
 */
export function getUserIdClient(): string | null {
  return Cookies.get("userId") || null;
}

/**
 * クライアントサイドでCookieにユーザーIDを設定する
 * @param userId - 設定するユーザーID
 */
export function setUserIdClient(userId: string): void {
  Cookies.set("userId", userId, { expires: 365 }); // 1年間有効
}
