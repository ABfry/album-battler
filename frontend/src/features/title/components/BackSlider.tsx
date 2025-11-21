"use client";
import React from "react";
import Image from "next/image";

import { Autoplay, EffectCoverflow } from "swiper/modules";
import { Swiper, SwiperSlide } from "swiper/react";
import "swiper/css";
import "swiper/css/pagination";

const imageUrls = [
  [
    "/title/sample.jpg",
    "/title/sample.jpg",
    "/title/sample.jpg",
    "/title/sample.jpg",
    "/title/sample.jpg",
    "/title/sample.jpg",
    "/title/sample.jpg",
    "/title/sample.jpg",
    "/title/sample.jpg",
  ],
  [
    "/title/sample.jpg",
    "/title/sample.jpg",
    "/title/sample.jpg",
    "/title/sample.jpg",
    "/title/sample.jpg",
    "/title/sample.jpg",
    "/title/sample.jpg",
    "/title/sample.jpg",
  ],
];

const slideSettings = {
  0: {
    slidesPerView: 1.4,
    spaceBetween: 10,
  },
  768: {
    slidesPerView: 2.2,
    spaceBetween: 10,
  },
  1024: {
    slidesPerView: 2.2,
    spaceBetween: 10,
  },
};

const reverseDirections = [false, true];

const TitleBackSlider = () => {
  return (
    <section className="fixed inset-0 -z-10">
      <div className="mx-auto flex h-full flex-col justify-center">
        {reverseDirections.map((reverseDirection, parentIndex) => {
          return (
            <Swiper
              key={parentIndex}
              modules={[Autoplay, EffectCoverflow]}
              breakpoints={slideSettings} // slidesPerViewを指定
              slidesPerView={"auto"} // ハイドレーションエラー対策
              centeredSlides={reverseDirection} // スライドを中央に配置
              loop={true} // スライドをループさせる
              autoplay={{
                delay: 0, // ディレイを0に設定
                disableOnInteraction: false, // ユーザーが操作しても自動再生を停止しない
                stopOnLastSlide: false, // 最後のスライドで停止しない
                waitForTransition: true, // トランジション待ち
                reverseDirection: reverseDirection,
              }}
              allowTouchMove={false}
              freeMode={{
                enabled: true,
                momentumRatio: 0.3,
                momentumVelocityRatio: 0.35,
              }}
              speed={8000}
              className="mx-auto my-4 block w-screen"
            >
              {imageUrls[parentIndex].map((imageUrl, index) => (
                <SwiperSlide key={index}>
                  <div className="relative aspect-[16/9] w-full">
                    <Image
                      src={imageUrl}
                      alt={`Gallery ${index + 1}`}
                      fill
                      className="rounded-2xl object-cover opacity-70"
                    />
                  </div>
                </SwiperSlide>
              ))}
            </Swiper>
          );
        })}
      </div>
    </section>
  );
};

export default TitleBackSlider;
