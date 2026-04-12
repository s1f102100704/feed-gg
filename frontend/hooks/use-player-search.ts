import { useCallback, useState } from "react";

import {
  ApiError,
  API_BASE_URL,
  Region,
  SearchResult,
} from "@/lib/player-search";

export function usePlayerSearch() {
  const [result, setResult] = useState<SearchResult | null>(null);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const searchPlayer = useCallback(
    async (region: Region, gameName: string, tagLine: string) => {
      setIsLoading(true);
      setError("");

      try {
        const response = await fetch(`${API_BASE_URL}/api/players/search`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            region,
            gameName,
            tagLine,
          }),
        });

        const payload = (await response.json()) as SearchResult | ApiError;

        if (!response.ok) {
          setResult(null);
          setError(
            "error" in payload ? payload.error : "プレイヤー情報の取得に失敗しました。",
          );
          return null;
        }

        const nextResult = payload as SearchResult;
        const normalizedResult: SearchResult = {
          ...nextResult,
          matches: nextResult.matches ?? [],
        };
        setResult(normalizedResult);
        return normalizedResult;
      } catch {
        setResult(null);
        setError("backend へ接続できませんでした。");
        return null;
      } finally {
        setIsLoading(false);
      }
    },
    [],
  );

  const clearResult = useCallback(() => {
    setResult(null);
    setError("");
  }, []);

  return {
    result,
    error,
    isLoading,
    searchPlayer,
    clearResult,
  };
}
