import { useState } from "react";

import {
  ApiError,
  API_BASE_URL,
  parseRiotID,
  Region,
  SearchResult,
} from "@/lib/player-search";

export function usePlayerSearch() {
  const [result, setResult] = useState<SearchResult | null>(null);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  async function searchPlayer(region: Region, riotId: string) {
    const parsed = parseRiotID(riotId);
    if (!parsed) {
      setError("Riot ID は `プレイヤー名#tagline` 形式で入力してください。");
      setResult(null);
      return;
    }

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
          gameName: parsed.gameName,
          tagLine: parsed.tagLine,
        }),
      });

      const payload = (await response.json()) as SearchResult | ApiError;

      if (!response.ok) {
        setResult(null);
        setError(
          "error" in payload ? payload.error : "プレイヤー情報の取得に失敗しました。",
        );
        return;
      }

      setResult(payload as SearchResult);
    } catch {
      setResult(null);
      setError("backend へ接続できませんでした。");
    } finally {
      setIsLoading(false);
    }
  }

  return {
    result,
    error,
    isLoading,
    searchPlayer,
  };
}
