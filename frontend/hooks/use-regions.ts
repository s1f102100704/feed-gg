"use client";

import { useEffect, useState } from "react";

import { API_BASE_URL } from "@/lib/player-search-api";
import { Region } from "@/types/player-search";

type RegionsResponse = {
  regions: string[];
};

export function useRegions() {
  const [regions, setRegions] = useState<Region[]>([]);
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;

    async function loadRegions() {
      try {
        const response = await fetch(`${API_BASE_URL}/api/regions`);
        const payload = (await response.json()) as RegionsResponse | { error?: string };

        if (!response.ok) {
          if (!cancelled) {
            setError(
              "error" in payload && payload.error
                ? payload.error
                : "region master の取得に失敗しました。",
            );
          }
          return;
        }

        if (!cancelled) {
          setRegions((payload as RegionsResponse).regions as Region[]);
          setError("");
        }
      } catch {
        if (!cancelled) {
          setError("region master の取得に失敗しました。");
        }
      }
    }

    void loadRegions();

    return () => {
      cancelled = true;
    };
  }, []);

  return {
    regions,
    error,
  };
}
