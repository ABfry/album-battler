import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/**
 * ファイルをBase64文字列に変換する
 * @param file 変換するファイル
 * @returns Base64文字列（Data URLプレフィックスなし）を返すPromise
 */
export function fileToBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();

    reader.onload = (event) => {
      const result = event.target?.result as string;
      if (!result) {
        reject(new Error("Failed to read file"));
        return;
      }

      // Data URLからBase64部分のみを抽出
      // 例: "data:image/png;base64,iVBORw0KG..." -> "iVBORw0KG..."
      const base64 = result.split(",")[1];
      if (!base64) {
        reject(new Error("Failed to extract base64 from data URL"));
        return;
      }

      resolve(base64);
    };

    reader.onerror = (error) => {
      reject(new Error(`FileReader error: ${error}`));
    };

    reader.readAsDataURL(file);
  });
}
