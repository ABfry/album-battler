import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

export function middleware(request: NextRequest) {
  const userId = request.cookies.get("userId")?.value;
  const { pathname } = request.nextUrl;

  // 除外するパス
  const excludedPaths = ["/", "/api-test", "/ui-demo", "/title"];

  // 静的ファイルとNext.js内部パスを除外
  if (
    pathname.startsWith("/api/") ||
    pathname.startsWith("/_next/") ||
    pathname.startsWith("/favicon.ico") ||
    pathname.match(/\.(svg|png|jpg|jpeg|gif|webp|ico)$/)
  ) {
    return NextResponse.next();
  }

  // 除外パスのチェック
  if (excludedPaths.includes(pathname)) {
    return NextResponse.next();
  }

  // userIdがない場合、ルートにリダイレクト
  if (!userId) {
    const url = request.nextUrl.clone();
    url.pathname = "/";
    return NextResponse.redirect(url);
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/((?!api|_next/static|_next/image|favicon.ico).*)"],
};
