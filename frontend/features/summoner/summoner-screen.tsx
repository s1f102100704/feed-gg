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
import { useLabels } from "@/hooks/use-labels";
import { usePlayerLabels } from "@/hooks/use-player-labels";
import { usePlayerSearch } from "@/hooks/use-player-search";
import { DEFAULT_REGION, Region } from "@/lib/regions";
import { buildSummonerPath } from "@/lib/player-search-path";

type SummonerScreenProps = {
  initialRegions: readonly Region[];
  resolvedRegion: Region | null;
  decodedGameName: string;
  decodedTagLine: string;
  initialRiotId: string;
};

export function SummonerScreen({
  initialRegions,
  resolvedRegion,
  decodedGameName,
  decodedTagLine,
  initialRiotId,
}: SummonerScreenProps) {
  const router = useRouter();
  const [region, setRegion] = useState<Region>(resolvedRegion ?? DEFAULT_REGION);
  const [riotId, setRiotId] = useState(initialRiotId);
  const { result, error, isLoading, searchPlayer, clearResult } = usePlayerSearch();
  const { labels, error: labelsError, isLoading: isLabelsLoading, fetchLabels } = useLabels();
  const {
    labels: playerLabels,
    labelsPUUID,
    totalVotes: totalLabelVotes,
    error: playerLabelsError,
    isLoading: isPlayerLabelsLoading,
    isSaving: isPlayerLabelSaving,
    fetchPlayerLabels,
    votePlayerLabel,
  } = usePlayerLabels();
  const resultPUUID = result?.puuid;
  const scopedPlayerLabels = labelsPUUID === resultPUUID ? playerLabels : [];
  const scopedTotalLabelVotes = labelsPUUID === resultPUUID ? totalLabelVotes : 0;
  const availableRegions = initialRegions.length > 0 ? initialRegions : [region];
  const selectedRegion = availableRegions.includes(region) ? region : availableRegions[0];

  useEffect(() => {
    void fetchLabels();
  }, [fetchLabels]);

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

  useEffect(() => {
    if (!resultPUUID) {
      return;
    }

    const controller = new AbortController();
    void fetchPlayerLabels(resultPUUID, controller.signal);

    return () => {
      controller.abort();
    };
  }, [fetchPlayerLabels, resultPUUID]);

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

        {error ? <SearchErrorMessage message={error} /> : null}

        {result ? (
          <div className="space-y-6">
            <SummonerHero
              key={result.puuid}
              availableLabels={labels}
              labelError={playerLabelsError || labelsError}
              playerLabels={scopedPlayerLabels}
              result={result}
              totalLabelVotes={scopedTotalLabelVotes}
              isLoadingPlayerLabels={isPlayerLabelsLoading || isLabelsLoading}
              isSavingPlayerLabel={isPlayerLabelSaving}
              onVotePlayerLabel={async (label) => {
                const response = await votePlayerLabel(result.puuid, label.id);
                return response !== null;
              }}
            />
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
