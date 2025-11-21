"use client";

import { useState } from "react";
import { battleApi } from "../lib/api/battleApi";
import { GetBattleResultResponse } from "../lib/api/types";

type UseResult = {
  getResult: (battleId: string) => Promise<GetBattleResultResponse | null>;
  loading: boolean;
  error: string | null;
};

export function useResult(): UseResult {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getResult = async (battleId: string) => {
    setLoading(true);
    setError(null);
    try {
      const response = await battleApi.getResult(battleId);
      const data: GetBattleResultResponse = {
        battle_id: response.battle_id,
        winner_user_id: response.winner_user_id,
        scores: response.scores,
      };
      return data;
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to fetch battle result";
      setError(message);
      return null;
    } finally {
      setLoading(false);
    }
  };

  return { getResult, loading, error };
}
