"use client";

import { FormEvent, useState } from "react";

import { FeedTitle } from "@/components/search/feed-title";
import { SearchForm } from "@/components/search/search-form";
import { SearchResultCard } from "@/components/search/search-result-card";
import { usePlayerSearch } from "@/hooks/use-player-search";
import { Region } from "@/lib/player-search";

export default function Home() {
  const [region, setRegion] = useState<Region>("JP1");
  const [riotId, setRiotId] = useState("");
  const { result, error, isLoading, searchPlayer } = usePlayerSearch();

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    await searchPlayer(region, riotId);
  }

  return (
    <main className="flex min-h-screen items-center justify-center bg-[#050816] px-6 text-white">
      <div className="w-full max-w-4xl">
        <div className="space-y-8 text-center">
          <FeedTitle />
          <SearchForm
            region={region}
            riotId={riotId}
            isLoading={isLoading}
            onRegionChange={setRegion}
            onRiotIDChange={setRiotId}
            onSubmit={handleSubmit}
          />

          {error ? <p className="text-sm text-red-300">{error}</p> : null}

          {result ? <SearchResultCard result={result} /> : null}
        </div>
      </div>
    </main>
  );
}
