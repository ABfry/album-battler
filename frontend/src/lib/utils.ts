import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import imageCompression from "browser-image-compression";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/**
 * ファイルサイズを人間が読みやすい形式に変換
 */
function formatFileSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
}

/**
 * Base64文字列のバイトサイズを計算
 */
function getBase64Size(base64: string): number {
  // Base64: 4文字 = 3バイト
  // パディング文字('=')を考慮
  const padding = (base64.match(/=/g) || []).length;
  return (base64.length * 3) / 4 - padding;
}

/**
 * FileをData URLに変換
 */
function fileToDataURL(file: File | Blob): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      const result = e.target?.result as string;
      resolve(result);
    };
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
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

/**
 * 画像を500KB以下に圧縮
 *
 * @param file 元画像ファイル
 * @param maxSizeMB 最大サイズ（デフォルト0.5MB = 500KB）
 * @returns 圧縮後のData URL
 *
 * @example
 * const dataUrl = await compressImage(file, 0.5);
 */
export async function compressImage(
  file: File,
  maxSizeMB: number = 0.5
): Promise<string> {
  console.log(
    `[imageCompression] 元ファイル: ${file.name}, サイズ: ${formatFileSize(file.size)}, タイプ: ${file.type}`
  );

  let processedFile: File | Blob = file;

  // ステップ1: HEIC/HEIF形式の場合、JPEGに変換
  if (file.type === "image/heic" || file.type === "image/heif") {
    console.log("[imageCompression] HEIC形式を検出、JPEG変換中...");
    try {
      // 動的インポート（ブラウザ専用ライブラリ）
      const heic2any = (await import("heic2any")).default;

      const convertedBlob = await heic2any({
        blob: file,
        toType: "image/jpeg",
        quality: 0.9,
      });

      // heic2anyは配列またはBlobを返す可能性がある
      const blob = Array.isArray(convertedBlob)
        ? convertedBlob[0]
        : convertedBlob;

      processedFile = new File([blob], file.name.replace(/\.heic$/i, ".jpg"), {
        type: "image/jpeg",
      });

      console.log(
        `[imageCompression] HEIC変換完了: ${formatFileSize(processedFile.size)}`
      );
    } catch (error) {
      console.error("[imageCompression] HEIC変換エラー:", error);
      throw new Error("HEIC画像の変換に失敗しました");
    }
  }

  // ステップ2: 画像圧縮（500KB以下に）
  try {
    // BlobをFileに変換（browser-image-compressionはFile型を要求）
    const fileToCompress =
      processedFile instanceof File
        ? processedFile
        : new File([processedFile], file.name, { type: processedFile.type });

    const options = {
      maxSizeMB: maxSizeMB,
      maxWidthOrHeight: 2048, // 最大解像度
      useWebWorker: true, // Web Workerで非同期処理
      fileType: "image/jpeg", // 出力形式をJPEGに統一
    };

    console.log("[imageCompression] 圧縮開始...");
    const compressedFile = await imageCompression(fileToCompress, options);

    console.log(
      `[imageCompression] 圧縮完了: ${formatFileSize(file.size)} → ${formatFileSize(compressedFile.size)}`
    );
    console.log(
      `[imageCompression] 圧縮率: ${((compressedFile.size / file.size) * 100).toFixed(1)}%`
    );

    // Data URLに変換
    const dataUrl = await fileToDataURL(compressedFile);

    // Base64部分のサイズを確認
    const base64 = dataUrl.split(",")[1];
    const base64Size = getBase64Size(base64);
    console.log(
      `[imageCompression] Base64サイズ: ${formatFileSize(base64Size)}`
    );

    if (base64Size > maxSizeMB * 1024 * 1024) {
      console.warn(
        `[imageCompression] 警告: Base64サイズが${maxSizeMB}MBを超えています`
      );
    }

    return dataUrl;
  } catch (error) {
    console.error("[imageCompression] 圧縮エラー:", error);
    throw new Error("画像の圧縮に失敗しました");
  }
}
