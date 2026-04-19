"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";

import { FeedTitle } from "@/components/search/feed-title";
import { SearchErrorMessage } from "@/components/search/search-error-message";
import { SearchForm } from "@/components/search/search-form";
import { DEFAULT_REGION, Region } from "@/lib/regions";
import { buildSummonerPath } from "@/lib/player-search-path";

type HomeScreenProps = {
  initialRegions: readonly Region[];
};

export function HomeScreen({ initialRegions }: HomeScreenProps) {
  const router = useRouter();
  const [region, setRegion] = useState<Region>(DEFAULT_REGION);
  const [riotId, setRiotId] = useState("");
  const [error, setError] = useState("");
  const availableRegions = initialRegions.length > 0 ? initialRegions : [region];
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

          {error ? <SearchErrorMessage message={error} centered /> : null}
        </div>
      </div>
    </main>
  );
}
