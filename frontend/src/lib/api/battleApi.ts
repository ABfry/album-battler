import { fetchApi } from "./client";
import type {
  CreateBattleRequest,
  CreateBattleResponse,
  GetBattleResponse,
  GetImageResponse,
  SendImageRequest,
  SendImageResponse,
} from "./types";

/**
 * バトル関連のAPI
 */
export const battleApi = {
  /**
   * バトルを作成（デバッグ用）
   * POST /battle
   */
  createBattle: async (roomId: string): Promise<CreateBattleResponse> => {
    return fetchApi<CreateBattleResponse>("/battle", {
      method: "POST",
      body: JSON.stringify({ room_id: roomId } satisfies CreateBattleRequest),
    });
  },

  /**
   * バトル情報を取得
   * GET /battle/{id}
   */
  getBattle: async (battleId: string): Promise<GetBattleResponse> => {
    return fetchApi<GetBattleResponse>(`/battle/${battleId}`);
  },

  /**
   * バトルの画像一覧を取得
   * GET /battle/{id}/image
   */
  getImages: async (battleId: string): Promise<GetImageResponse> => {
    return fetchApi<GetImageResponse>(`/battle/${battleId}/image`);
  },

  /**
   * 画像を送信
   * POST /battle/{id}/send-image
   */
  sendImage: async (
    battleId: string,
    userId: string,
    imageBase64: string
  ): Promise<SendImageResponse> => {
    return fetchApi<SendImageResponse>(`/battle/${battleId}/send-image`, {
      method: "POST",
      body: JSON.stringify({
        user_id: userId,
        image_base64: imageBase64,
      } satisfies SendImageRequest),
    });
  },
};
