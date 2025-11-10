import { RoomPage } from "@/src/features/room/containers/RoomPage";

export default function Page() {
  // 「ルーティングとデータ取得の起点となるServer Component」として最小責務だけを持つ
  return <RoomPage />;
}
