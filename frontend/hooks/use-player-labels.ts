import { useCallback, useState } from "react";

import { API_BASE_URL } from "@/lib/player-search-api";
import {
  ApiError,
  PlayerLabelsResponse,
  PlayerLabelVoteResponse,
} from "@/types/player-search";

function isAbortError(error: unknown) {
  return error instanceof DOMException && error.name === "AbortError";
}

export function usePlayerLabels() {
  const [labels, setLabels] = useState<PlayerLabelsResponse>({
    labels: [],
    totalVotes: 0,
  });
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  const fetchPlayerLabels = useCallback(async (puuid: string, signal?: AbortSignal) => {
    setIsLoading(true);
    setError("");

    try {
      const response = await fetch(
        `${API_BASE_URL}/api/players/${encodeURIComponent(puuid)}/labels`,
        { signal },
      );
      const payload = (await response.json()) as PlayerLabelsResponse | ApiError;

      if (signal?.aborted) {
        return null;
      }

      if (!response.ok) {
        const nextError =
          "error" in payload ? payload.error : "プレイヤーのラベル取得に失敗しました。";
        setLabels({ labels: [], totalVotes: 0 });
        setError(nextError);
        return null;
      }

      const nextLabels = payload as PlayerLabelsResponse;
      setLabels({
        labels: nextLabels.labels ?? [],
        totalVotes: nextLabels.totalVotes ?? 0,
      });
      return nextLabels;
    } catch (error) {
      if (isAbortError(error) || signal?.aborted) {
        return null;
      }

      setLabels({ labels: [], totalVotes: 0 });
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
      const response = await fetch(
        `${API_BASE_URL}/api/players/${encodeURIComponent(puuid)}/labels`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ labelId }),
        },
      );
      const payload = (await response.json()) as PlayerLabelVoteResponse | ApiError;

      if (!response.ok) {
        setError("error" in payload ? payload.error : "ラベルの保存に失敗しました。");
        return null;
      }

      const nextLabels = payload as PlayerLabelVoteResponse;
      setLabels({
        labels: nextLabels.labels ?? [],
        totalVotes: nextLabels.totalVotes ?? 0,
      });
      return nextLabels;
    } catch {
      setError("ラベルの保存に失敗しました。");
      return null;
    } finally {
      setIsSaving(false);
    }
  }, []);

  return {
    labels: labels.labels,
    totalVotes: labels.totalVotes,
    error,
    isLoading,
    isSaving,
    fetchPlayerLabels,
    votePlayerLabel,
  };
}
