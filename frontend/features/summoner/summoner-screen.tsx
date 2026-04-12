"use client";

import { FormEvent, useEffect, useState } from "react";
import { useRouter } from "next/navigation";

import { MatchHistoryPanel } from "@/components/match/match-history-panel";
import { RankPanel } from "@/components/rank/rank-panel";
import { SearchErrorMessage } from "@/components/search/search-error-message";
import { SearchForm } from "@/components/search/search-form";
import { SearchLoadingState } from "@/components/search/search-loading-state";
import { SummonerHero } from "@/components/search/summoner-hero";
import { SummonerContentLayout } from "@/components/summoner/summoner-content-layout";
import { usePlayerSearch } from "@/hooks/use-player-search";
import { useRegions } from "@/hooks/use-regions";
import { buildSummonerPath } from "@/lib/player-search-path";
import { Region } from "@/types/player-search";

type SummonerScreenProps = {
  resolvedRegion: Region | null;
  decodedGameName: string;
  decodedTagLine: string;
  initialRiotId: string;
};

export function SummonerScreen({
  resolvedRegion,
  decodedGameName,
  decodedTagLine,
  initialRiotId,
}: SummonerScreenProps) {
  const router = useRouter();
  const [region, setRegion] = useState<Region>(resolvedRegion ?? "JP1");
  const [riotId, setRiotId] = useState(initialRiotId);
  const { result, error, isLoading, searchPlayer, clearResult } = usePlayerSearch();
  const { regions, error: regionsError } = useRegions();
  const availableRegions = regions.length > 0 ? regions : [region];
  const selectedRegion = availableRegions.includes(region) ? region : availableRegions[0];

  useEffect(() => {
    if (!resolvedRegion) {
      clearResult();
      return;
    }

    void searchPlayer(resolvedRegion, decodedGameName, decodedTagLine);
  }, [
    clearResult,
    decodedGameName,
    decodedTagLine,
    resolvedRegion,
    searchPlayer,
  ]);

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const nextPath = buildSummonerPath(selectedRegion, riotId);
    if (!nextPath) {
      return;
    }

    router.push(nextPath);
  }

  return (
    <main className="min-h-screen bg-[#050816] px-6 py-8 text-white">
      <div className="mx-auto w-full max-w-6xl space-y-8">
        <div className="space-y-5">
          <div>
            <p className="text-sm uppercase tracking-[0.35em] text-cyan-300">
              feed.gg
            </p>
          </div>

          <SearchForm
            regions={availableRegions}
            region={selectedRegion}
            riotId={riotId}
            isLoading={isLoading}
            onRegionChange={setRegion}
            onRiotIDChange={setRiotId}
            onSubmit={handleSubmit}
          />
        </div>

        {error || regionsError ? (
          <SearchErrorMessage message={error || regionsError} />
        ) : null}

        {result ? (
          <div className="space-y-6">
            <SummonerHero result={result} />
            <SummonerContentLayout
              left={<RankPanel result={result} />}
              right={<MatchHistoryPanel result={result} />}
            />
          </div>
        ) : (
          <SearchLoadingState
            message={
              resolvedRegion
                ? "プレイヤー情報を読み込んでいます。"
                : "region が不正です。検索バーからやり直してください。"
            }
          />
        )}
      </div>
    </main>
  );
}
