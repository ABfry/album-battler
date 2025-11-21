"use client";

import { useParams } from "next/navigation";
import { ResultPage } from "@/src/features/result/containers/ResultPage";

export default function Page() {
  const { battleId } = useParams() as { battleId: string };

  return <ResultPage battleId={battleId} />;
}
