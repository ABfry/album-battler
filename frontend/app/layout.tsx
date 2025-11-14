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

const wsUrl = process.env.NEXT_PUBLIC_WEBSOCKET_URL || "ws://localhost:8080/ws";

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
