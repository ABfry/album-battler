import type { Metadata } from "next";
import {
  Geist,
  Geist_Mono,
  Noto_Sans_JP,
  M_PLUS_Rounded_1c,
} from "next/font/google";
import "./globals.css";
import { ClientLayout } from "./ClientLayout";
import { getUserId } from "@/src/lib/auth/getUserId";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

const notoSansJP = Noto_Sans_JP({
  variable: "--font-noto-sans-jp",
  subsets: ["latin"],
  weight: ["400", "700", "900"],
});

const mPlusRounded1c = M_PLUS_Rounded_1c({
  variable: "--font-m-plus-rounded-1c",
  subsets: ["latin"],
  weight: ["400"],
});

export const metadata: Metadata = {
  title: "アルバムバトラー",
  description: "アルバムで対戦するゲーム",
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  // サーバーサイドでCookieからユーザーIDを取得
  const userId = await getUserId();

  return (
    <html lang="ja">
      <body
        className={`${geistSans.variable} ${geistMono.variable} ${notoSansJP.variable} ${mPlusRounded1c.variable} font-m-plus-rounded-1c antialiased`}
        style={{
          backgroundImage: "url(/background.png)",
          backgroundRepeat: "repeat",
          backgroundSize: "auto",
        }}
      >
        <ClientLayout userId={userId}>{children}</ClientLayout>
      </body>
    </html>
  );
}
