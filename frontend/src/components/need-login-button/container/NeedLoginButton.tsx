"use client";

import { useState, useEffect } from "react";
import { NeedLoginButtonView } from "../components/NeedLoginButtonView";
import { createUser } from "../api/createUser";

export type NeedLoginButtonProps = {
  onClick?: () => void; // ログインの有無に関わらずクリックした時の処理
  loggedInOnClick?: () => void; // ログインしている時のクリック処理
  notLoggedInOnClick?: () => void; // ログインしていない時のクリック処理
  onLoginSuccess?: () => void; // ユーザー作成成功時のコールバック
  content: string;
  loggedInLink?: string; // ログインしている時、クリック時のリンク先
  className?: string;
};

export function NeedLoginButton({
  onClick,
  loggedInOnClick,
  notLoggedInOnClick,
  onLoginSuccess,
  content,
  loggedInLink,
  className,
}: NeedLoginButtonProps) {
  const [userId, setUserId] = useState<string | null>(null);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [userName, setUserName] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // クッキーからuserIDを取得
  useEffect(() => {
    if (typeof window !== "undefined") {
      const cookieUserId = document.cookie
        .split("; ")
        .find((row) => row.startsWith("userId="))
        ?.split("=")[1];
      setUserId(cookieUserId || null);
    }
  }, []);

  const handleButtonClick = () => {
    onClick?.();
    if (userId) {
      loggedInOnClick?.();
    } else {
      notLoggedInOnClick?.();
      setIsDialogOpen(true);
    }
  };

  const handleSubmit = async () => {
    if (!userName.trim()) {
      setError("ユーザー名を入力してください");
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const newUserId = await createUser(userName);
      if (typeof window !== "undefined") {
        // クッキーに保存（10年間有効 = 実質無期限）
        const expires = new Date();
        expires.setFullYear(expires.getFullYear() + 10);
        document.cookie = `userId=${newUserId}; expires=${expires.toUTCString()}; path=/`;
      }
      setUserId(newUserId);
      setIsDialogOpen(false);
      setUserName("");
      onLoginSuccess?.();
      onClick?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : "エラーが発生しました");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <NeedLoginButtonView
      isLoggedIn={!!userId}
      content={content}
      loggedInLink={loggedInLink}
      className={className}
      isDialogOpen={isDialogOpen}
      userName={userName}
      isLoading={isLoading}
      error={error}
      onButtonClick={handleButtonClick}
      onDialogOpenChange={setIsDialogOpen}
      onUserNameChange={setUserName}
      onSubmit={handleSubmit}
    />
  );
}
