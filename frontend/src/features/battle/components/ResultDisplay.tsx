import type { GetBattleResultResponse, UserInfo } from "@/src/lib/api/types";
import Image from "next/image";
import { animate, motion } from "framer-motion";
import { useEffect, useRef, useState } from "react";

type ResultDisplayProps = {
  battleResult: GetBattleResultResponse;
  players: UserInfo[];
};

/**
 * バトル結果表示コンポーネント (Presentational)
 */
export function ResultDisplay({ battleResult, players }: ResultDisplayProps) {
  const winnerPlayer = players.find(
    (p) => p.id === battleResult.winner_user_id
  );

  // アニメーション設定
  const maxFinalScore = Math.max(
    1,
    ...battleResult.results.map((r) => r.final_score)
  );
  const MAX_BAR_HEIGHT = 300;
  const AI_BAR_DURATION = 4.0;
  const CLAP_BAR_DURATION = 0.5;
  const PAUSE_DURATION = 1.5;
  const BAR_DELAY_STEP = 0.15;

  const lastIndexDelay = Math.max(battleResult.results.length - 1, 0);
  const totalBarTimeline =
    AI_BAR_DURATION +
    PAUSE_DURATION +
    CLAP_BAR_DURATION +
    lastIndexDelay * BAR_DELAY_STEP;
  const winnerStart = totalBarTimeline + 0.2;

  const drumRollRef = useRef<HTMLAudioElement | null>(null);
  const showResultSoundRef = useRef<HTMLAudioElement | null>(null);

  // 結果演出にあわせたSE制御
  useEffect(() => {
    if (typeof window === "undefined") return;

    const drum = new Audio("/sounds/drum-roll.mp3");
    drum.loop = true;
    drum.volume = 0.7;
    drum.preload = "auto";
    drumRollRef.current = drum;

    const show = new Audio("/sounds/show-result-sound.mp3");
    show.volume = 0.8;
    show.preload = "auto";
    showResultSoundRef.current = show;

    const aiPhaseDuration = AI_BAR_DURATION + lastIndexDelay * BAR_DELAY_STEP;
    void drum.play();

    const stopDrumTimer = window.setTimeout(() => {
      drum.loop = false;
      drum.pause();
      drum.currentTime = 0;
    }, Math.max(0, (aiPhaseDuration + 0.3) * 1000));

    const hasClapScore = battleResult.results.some(
      (r) => Math.max(r.user_score, 0) > 0
    );
    const clapStartDelay = AI_BAR_DURATION + PAUSE_DURATION;
    const showResultTimer = hasClapScore
      ? window.setTimeout(() => {
          if (!showResultSoundRef.current) return;
          try {
            showResultSoundRef.current.currentTime = 0;
            void showResultSoundRef.current.play();
          } catch {
            // 自動再生制限などは無視
          }
        }, Math.max(0, clapStartDelay * 1000))
      : null;

    return () => {
      clearTimeout(stopDrumTimer);
      if (showResultTimer) clearTimeout(showResultTimer);
      drum.pause();
      drum.currentTime = 0;
      show.pause();
      show.currentTime = 0;
    };
  }, [battleResult.results, lastIndexDelay]);

  // プレイヤーごとの色定義（順番に割り当て）
  const playerColors = [
    {
      aiBar: "bg-gradient-to-t from-red-500 to-red-400",
      clapBar: "bg-red-200",
      text: "text-red-600",
    },
    {
      aiBar: "bg-gradient-to-t from-blue-500 to-blue-400",
      clapBar: "bg-blue-200",
      text: "text-blue-600",
    },
    {
      aiBar: "bg-gradient-to-t from-yellow-500 to-yellow-400",
      clapBar: "bg-yellow-200",
      text: "text-yellow-600",
    },
    {
      aiBar: "bg-gradient-to-t from-green-500 to-green-400",
      clapBar: "bg-green-200",
      text: "text-green-600",
    },
    {
      aiBar: "bg-gradient-to-t from-orange-500 to-orange-400",
      clapBar: "bg-orange-200",
      text: "text-orange-600",
    },
  ];

  return (
    <div className="flex min-h-screen flex-col items-center justify-center p-4">
      <div className="relative flex h-screen w-full max-w-4xl flex-col items-center justify-center py-8">
        {/* 勝者発表（中央） */}
        <motion.div
          initial={{ opacity: 0, scale: 0 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{
            delay: winnerStart + 0.1,
            duration: 0.8,
          }}
          className="mb-16 space-y-4 text-center"
        >
          <div className="flex justify-center">
            <Image
              src="/winner-text.png"
              alt="Winner"
              width={240}
              height={80}
              priority
            />
          </div>
          <div className="text-4xl font-black text-slate-800">
            {winnerPlayer?.name || "Unknown Player"}
          </div>
        </motion.div>

        {/* プレイヤー一覧と棒グラフ（下部） */}
        <div className="w-full">
          {/* 点数を横一列で固定表示 */}
          <div className="mb-2 flex justify-center gap-4 md:gap-10">
            {battleResult.results.map((result, index) => {
              const color = playerColors[index % playerColors.length];
              return (
                <div
                  key={`score-${result.user_id}`}
                  className="flex w-16 flex-col items-center text-center md:w-20"
                >
                  <AnimatedScore
                    aiScore={Math.max(result.ai_score, 0)}
                    clapScore={Math.max(result.user_score, 0)}
                    delay={index * BAR_DELAY_STEP}
                    aiDuration={AI_BAR_DURATION}
                    pauseDuration={
                      Math.max(result.user_score, 0) > 0 ? PAUSE_DURATION : 0
                    }
                    clapDuration={
                      Math.max(result.user_score, 0) > 0 ? CLAP_BAR_DURATION : 0
                    }
                    className={`text-2xl font-black ${color.text} whitespace-nowrap`}
                  />
                  <div className="mt-4" />
                </div>
              );
            })}
          </div>

          {/* 棒グラフ */}
          <div className="mb-4 flex items-end justify-center gap-4 md:gap-10">
            {battleResult.results.map((result, index) => {
              const player = players.find((p) => p.id === result.user_id);
              const isWinner = result.user_id === battleResult.winner_user_id;
              const finalScore = Math.max(result.final_score, 0);
              const aiScore = Math.max(result.ai_score, 0);
              const barHeight = (finalScore / maxFinalScore) * MAX_BAR_HEIGHT;
              const aiRatio = finalScore > 0 ? aiScore / finalScore : 0;
              const aiHeight = barHeight * aiRatio;
              const clapHeight = Math.max(barHeight - aiHeight, 0);
              const hasClap = Math.max(result.user_score, 0) > 0;
              const color = playerColors[index % playerColors.length];

              return (
                <div
                  key={result.user_id}
                  className="flex flex-col items-center"
                >
                  <div
                    className="relative w-16 md:w-20"
                    style={{ height: Math.max(barHeight, 20) }}
                  >
                    <motion.div
                      className={`absolute right-0 bottom-0 left-0 ${color.aiBar}`}
                      initial={{ height: 0 }}
                      animate={{ height: aiHeight }}
                      transition={{
                        delay: index * BAR_DELAY_STEP,
                        duration: AI_BAR_DURATION,
                        ease: "easeOut",
                      }}
                    />
                    <motion.div
                      className={`absolute right-0 bottom-0 left-0 ${color.clapBar}`}
                      initial={{ height: 0 }}
                      animate={{ height: clapHeight }}
                      transition={{
                        delay:
                          index * BAR_DELAY_STEP +
                          AI_BAR_DURATION +
                          (hasClap ? PAUSE_DURATION : 0),
                        duration: hasClap ? CLAP_BAR_DURATION : 0,
                        ease: "easeOut",
                      }}
                      style={{ bottom: aiHeight }}
                    />
                    {isWinner && (
                      <div className="absolute top-6 left-1/2 flex -translate-x-1/2 items-center justify-center">
                        <motion.div
                          className="absolute inset-0 flex items-center justify-center"
                          initial={{ opacity: 0, scale: 0, rotate: 0 }}
                          animate={{ opacity: 0.9, scale: 2, rotate: 360 }}
                          transition={{
                            delay:
                              index * BAR_DELAY_STEP +
                              AI_BAR_DURATION +
                              (hasClap ? PAUSE_DURATION : 0) +
                              (hasClap ? CLAP_BAR_DURATION : 0) +
                              0.1,
                            duration: 0.4,
                            ease: "easeOut",
                            rotate: {
                              delay:
                                index * BAR_DELAY_STEP +
                                AI_BAR_DURATION +
                                (hasClap ? PAUSE_DURATION : 0) +
                                (hasClap ? CLAP_BAR_DURATION : 0) +
                                0.1,
                              duration: 2.4,
                              repeat: Infinity,
                              ease: "linear",
                            },
                          }}
                        >
                          <Image
                            src="/back-light.png"
                            alt="back light"
                            width={72}
                            height={72}
                            priority
                          />
                        </motion.div>
                        <motion.div
                          className="relative flex items-center justify-center"
                          initial={{ opacity: 0, scale: 0 }}
                          animate={{ opacity: 1, scale: 1.25 }}
                          transition={{
                            delay:
                              index * BAR_DELAY_STEP +
                              AI_BAR_DURATION +
                              (hasClap ? PAUSE_DURATION : 0) +
                              (hasClap ? CLAP_BAR_DURATION : 0) +
                              0.1,
                            duration: 0.3,
                          }}
                        >
                          <Image
                            src="/crown-icon.png"
                            alt="crown"
                            width={36}
                            height={36}
                            priority
                          />
                        </motion.div>
                      </div>
                    )}
                  </div>

                  <div className="mt-4 h-16 w-16 overflow-hidden rounded-full border-4 border-slate-300 bg-slate-200 shadow-lg md:h-20 md:w-20">
                    <Image
                      src={player?.icon_url || ""}
                      alt={player?.name || "Unknown"}
                      width={80}
                      height={80}
                      className="h-full w-full object-cover"
                      unoptimized
                    />
                  </div>

                  <div className="mt-2 w-16 text-center md:w-20">
                    <p className="truncate text-sm font-bold text-slate-700 md:text-base">
                      {player?.name || "Unknown"}
                    </p>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}

type AnimatedScoreProps = {
  aiScore: number;
  clapScore: number;
  delay: number;
  aiDuration: number;
  pauseDuration: number;
  clapDuration: number;
  className?: string;
};

function AnimatedScore({
  aiScore,
  clapScore,
  delay,
  aiDuration,
  pauseDuration,
  clapDuration,
  className,
}: AnimatedScoreProps) {
  const [display, setDisplay] = useState(0);

  useEffect(() => {
    const totalScore = Math.round(aiScore + clapScore);
    const aiTarget = Math.round(aiScore);

    const aiControls = animate(0, aiTarget, {
      delay,
      duration: aiDuration,
      ease: "easeOut",
      onUpdate: (v) => setDisplay(Math.round(v)),
    });

    const clapTimer = setTimeout(
      () => {
        animate(aiTarget, totalScore, {
          duration: clapDuration,
          ease: "easeOut",
          onUpdate: (v) => setDisplay(Math.round(v)),
        });
      },
      (delay + aiDuration + pauseDuration) * 1000
    );

    return () => {
      aiControls.stop();
      clearTimeout(clapTimer);
    };
  }, [aiScore, clapScore, delay, aiDuration, pauseDuration, clapDuration]);

  return <div className={className}>{display}点</div>;
}
