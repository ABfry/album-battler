import type { Metadata } from "next";
import { Geist, Geist_Mono, Noto_Sans_JP } from "next/font/google";
import "./globals.css";
import { WebSocketProvider } from "@/src/lib/websocket/providers/WebSocketProvider";

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

export const metadata: Metadata = {
  title: "アルバムバトラー",
  description: "アルバムで対戦するゲーム",
};

// TODO: 本来はログインユーザーIDを使うべき（現在はハードコード）
const TEST_USER_ID = "550e8400-e29b-41d4-a716-446655440001";
const wsUrl = `${process.env.NEXT_PUBLIC_WEBSOCKET_URL || "ws://localhost:8080/ws"}?user_id=${TEST_USER_ID}`;

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja">
      <body
        className={`${geistSans.variable} ${geistMono.variable} ${notoSansJP.variable} antialiased`}
        style={{
          backgroundImage: "url(/background.png)",
          backgroundRepeat: "repeat",
          backgroundSize: "auto",
        }}
      >
        <WebSocketProvider url={wsUrl}>{children}</WebSocketProvider>
      </body>
    </html>
  );
}
