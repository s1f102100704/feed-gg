import { useCallback, useState } from "react";

import {
  fetchPlayerLabels as fetchPlayerLabelsAPI,
  votePlayerLabel as votePlayerLabelAPI,
} from "@/lib/player-labels-api";
import { PlayerLabelsResponse } from "@/types/player-labels";

type PlayerLabelsState = PlayerLabelsResponse & {
  puuid: string;
};

function isAbortError(error: unknown) {
  return error instanceof DOMException && error.name === "AbortError";
}

export function usePlayerLabels() {
  const [labels, setLabels] = useState<PlayerLabelsState>({
    puuid: "",
    labels: [],
    totalVotes: 0,
  });
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  const fetchPlayerLabels = useCallback(async (puuid: string, signal?: AbortSignal) => {
    setIsLoading(true);
    setError("");
    setLabels({ puuid, labels: [], totalVotes: 0 });

    try {
      const payload = await fetchPlayerLabelsAPI(puuid, signal);

      if (signal?.aborted) {
        return null;
      }

      if ("error" in payload) {
        setLabels({ puuid, labels: [], totalVotes: 0 });
        setError(payload.error);
        return null;
      }

      const nextLabels = payload as PlayerLabelsResponse;
      setLabels({
        puuid,
        labels: nextLabels.labels ?? [],
        totalVotes: nextLabels.totalVotes ?? 0,
      });
      return nextLabels;
    } catch (error) {
      if (isAbortError(error) || signal?.aborted) {
        return null;
      }

      setLabels({ puuid, labels: [], totalVotes: 0 });
      setError("プレイヤーのラベル取得に失敗しました。");
      return null;
    } finally {
      if (!signal?.aborted) {
        setIsLoading(false);
      }
    }
  }, []);

  const votePlayerLabel = useCallback(async (puuid: string, labelId: number) => {
    setIsSaving(true);
    setError("");

    try {
      const payload = await votePlayerLabelAPI(puuid, labelId);

      if ("error" in payload) {
        setError(payload.error);
        return null;
      }

      setLabels({
        puuid,
        labels: payload.labels ?? [],
        totalVotes: payload.totalVotes ?? 0,
      });
      return payload;
    } catch {
      setError("ラベルの保存に失敗しました。");
      return null;
    } finally {
      setIsSaving(false);
    }
  }, []);

  return {
    labels: labels.labels,
    labelsPUUID: labels.puuid,
    totalVotes: labels.totalVotes,
    error,
    isLoading,
    isSaving,
    fetchPlayerLabels,
    votePlayerLabel,
  };
}
