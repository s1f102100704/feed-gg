"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";

import { FeedTitle } from "@/components/search/feed-title";
import { SearchErrorMessage } from "@/components/search/search-error-message";
import { SearchForm } from "@/components/search/search-form";
import { buildSummonerPath, Region } from "@/lib/player-search";

export function HomeScreen() {
  const router = useRouter();
  const [region, setRegion] = useState<Region>("JP1");
  const [riotId, setRiotId] = useState("");
  const [error, setError] = useState("");

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const nextPath = buildSummonerPath(region, riotId);

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
            region={region}
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
