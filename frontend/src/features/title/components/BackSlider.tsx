"use client";
import { useState, useEffect } from "react";
import Image from "next/image";

const IMAGE_COUNT = 21;
interface FallingImage {
  id: number;
  left: number;
  duration: number;
  delay: number;
  size: number;
  imageUrl: string;
}

/**
 * タイトル画面の背景スライダー
 */
const TitleBackSlider = () => {
  const [mounted, setMounted] = useState(false);
  const [fallingImages, setFallingImages] = useState<FallingImage[]>([]);

  // クライアント側でのみランダム値を生成
  useEffect(() => {
    const images = Array.from({ length: 10 }, (_, i) => {
      const randomImageNum = Math.floor(Math.random() * IMAGE_COUNT) + 1;
      return {
        id: i,
        left: Math.random() * 100,
        duration: 10 + Math.random() * 10,
        delay: Math.random() * -15,
        size: 150 + Math.random() * 150,
        imageUrl: `/title/img${randomImageNum}.jpg`,
      };
    });

    // eslint-disable-next-line react-hooks/set-state-in-effect
    setFallingImages(images);
    setMounted(true);
  }, []);

  if (!mounted) {
    return null;
  }

  return (
    <>
      <style jsx>{`
        @keyframes fall {
          from {
            transform: translateY(-100%) var(--rotation);
            opacity: 0;
          }
          10% {
            opacity: 1;
          }
          90% {
            opacity: 1;
          }
          to {
            transform: translateY(calc(100vh + 100%)) var(--rotation);
            opacity: 0;
          }
        }
      `}</style>
      <section
        className="fixed inset-0 -z-10 overflow-hidden"
        style={{ perspective: "1000px" }}
      >
        {fallingImages.map((img) => {
          // 画面中心（50%）からの距離に応じて回転角度を計算
          const distanceFromCenter = img.left - 50;
          const rotateY = distanceFromCenter * -0.5;
          const rotateX = 15;

          return (
            <div
              key={img.id}
              className="absolute"
              style={{
                left: `${img.left}%`,
                width: `${img.size}px`,
                height: `${img.size * 0.75}px`, // 4:3のアスペクト比
                animation: `fall ${img.duration}s linear ${img.delay}s infinite`,
                transformStyle: "preserve-3d",
                // @ts-expect-error - CSS変数のため
                "--rotation": `rotateY(${rotateY}deg) rotateX(${rotateX}deg)`,
              }}
            >
              <Image
                src={img.imageUrl}
                alt=""
                fill
                className="border-2 border-white object-cover"
              />
            </div>
          );
        })}
      </section>
    </>
  );
};

export default TitleBackSlider;
