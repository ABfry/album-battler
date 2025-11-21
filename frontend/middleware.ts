import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

export function middleware(request: NextRequest) {
  const userId = request.cookies.get("userId")?.value;
  const { pathname } = request.nextUrl;

  // 認証不要なパス (動的な除外リスト)
  const excludedPaths = ["/", "/api-test", "/ui-demo", "/title"];

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

// Matcher: 静的リソースを完全除外 (API, _next, 画像ファイル)
// Middleware: 認証が必要なページの動的制御
export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * - Static files (svg, png, jpg, jpeg, gif, webp, ico, json)
     */
    "/((?!api|_next/static|_next/image|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|gif|webp|ico|json)$).*)",
  ],
};
