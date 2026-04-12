"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";

import { FeedTitle } from "@/components/search/feed-title";
import { SearchErrorMessage } from "@/components/search/search-error-message";
import { SearchForm } from "@/components/search/search-form";
import { useRegions } from "@/hooks/use-regions";
import { buildSummonerPath } from "@/lib/player-search-path";
import { Region } from "@/types/player-search";

export function HomeScreen() {
  const router = useRouter();
  const [region, setRegion] = useState<Region>("JP1");
  const [riotId, setRiotId] = useState("");
  const [error, setError] = useState("");
  const { regions, error: regionsError } = useRegions();
  const availableRegions = regions.length > 0 ? regions : [region];
  const selectedRegion = availableRegions.includes(region) ? region : availableRegions[0];

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const nextPath = buildSummonerPath(selectedRegion, riotId);

    if (!nextPath) {
      setError("Riot ID は `プレイヤー名#tagline` 形式で入力してください。");
      return;
    }

    setError("");
    router.push(nextPath);
  }

  return (
    <main className="flex min-h-screen items-center justify-center bg-[#050816] px-6 text-white">
      <div className="w-full max-w-4xl">
        <div className="space-y-8 text-center">
          <FeedTitle />
          <SearchForm
            regions={availableRegions}
            region={selectedRegion}
            riotId={riotId}
            isLoading={false}
            onRegionChange={setRegion}
            onRiotIDChange={setRiotId}
            onSubmit={handleSubmit}
          />

          {error || regionsError ? (
            <SearchErrorMessage message={error || regionsError} centered />
          ) : null}
        </div>
      </div>
    </main>
  );
}
