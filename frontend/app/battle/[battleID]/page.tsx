import { BattlePage } from "@/src/features/battle/containers/BattlePage";

export default async function Page({
  params,
}: {
  params: Promise<{ battleID: string }>;
}) {
  const { battleID } = await params;

  return <BattlePage battleID={battleID} />;
}
