import { useCallback, useState } from "react";

import { API_BASE_URL } from "@/lib/player-search-api";
import { ApiError, Label } from "@/types/player-search";

type LabelsResponse = {
  labels: Label[];
};

export function useLabels() {
  const [labels, setLabels] = useState<Label[]>([]);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const fetchLabels = useCallback(async () => {
    setIsLoading(true);
    setError("");

    try {
      const response = await fetch(`${API_BASE_URL}/api/labels`);
      const payload = (await response.json()) as LabelsResponse | ApiError;

      if (!response.ok) {
        setLabels([]);
        setError(
          "error" in payload ? payload.error : "ラベル一覧の取得に失敗しました。",
        );
        return [];
      }

      const nextLabels = (payload as LabelsResponse).labels ?? [];
      setLabels(nextLabels);
      return nextLabels;
    } catch {
      setLabels([]);
      setError("ラベル一覧の取得に失敗しました。");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  return {
    labels,
    error,
    isLoading,
    fetchLabels,
  };
}
