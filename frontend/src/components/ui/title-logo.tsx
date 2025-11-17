import Image from "next/image";
import { cn } from "@/src/lib/utils";

const ASPECT_RATIO_VALUE = `${963} / ${486}` as const;

export interface TitleLogoProps {
  /**
   * 代替テキスト
   * @default "Album Battler Logo"
   */
  alt?: string;
  /**
   * CSSクラス（必須）
   * w-*, h-*でサイズ指定（例: w-full, h-full, w-64）
   */
  className: string;
  /**
   * 画像の読み込み方法
   * @default "lazy"
   */
  loading?: "lazy" | "eager";
}

/**
 * classNameに指定されたTailwindクラスからwidth/heightのクラスを検出
 */
function hasWidthClass(className: string): boolean {
  return /\bw-/.test(className);
}

function hasHeightClass(className: string): boolean {
  return /\bh-/.test(className);
}

/**
 * Album Battlerのタイトルロゴコンポーネント
 * アスペクト比 963:486 で固定されています
 *
 * @example
 * // 幅指定
 * <TitleLogo className="w-full" />
 * <TitleLogo className="w-64" />
 *
 * @example
 * // 高さ指定
 * <TitleLogo className="h-full" />
 * <TitleLogo className="h-32" />
 */
export function TitleLogo({
  alt = "Album Battler Logo",
  className,
  loading = "lazy",
}: TitleLogoProps) {
  const hasWidthInClassName = hasWidthClass(className);
  const hasHeightInClassName = hasHeightClass(className);

  // バリデーション: classNameにw-*またはh-*が必要
  if (!hasWidthInClassName && !hasHeightInClassName) {
    throw new Error(
      "TitleLogo: className に w-* または h-* を指定してください（例: w-full, h-64）"
    );
  }

  return (
    <div
      className={cn("relative", className)}
      style={{ aspectRatio: ASPECT_RATIO_VALUE }}
    >
      <Image
        src="/album-battler-logo.png"
        alt={alt}
        fill
        loading={loading}
        className="object-contain"
        priority={loading === "eager"}
      />
    </div>
  );
}
