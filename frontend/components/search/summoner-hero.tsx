"use client";

import Image from "next/image";
import { useMemo, useState } from "react";

import { Label, PlayerLabelSummary } from "@/types/player-labels";
import { SearchResult } from "@/types/player-search";

type SummonerHeroProps = {
  availableLabels: Label[];
  labelError: string;
  playerLabels: PlayerLabelSummary[];
  result: SearchResult;
  totalLabelVotes: number;
  isLoadingPlayerLabels: boolean;
  isSavingPlayerLabel: boolean;
  onVotePlayerLabel: (label: Label) => Promise<boolean>;
};

export function SummonerHero({
  availableLabels,
  labelError,
  playerLabels,
  result,
  totalLabelVotes,
  isLoadingPlayerLabels,
  isSavingPlayerLabel,
  onVotePlayerLabel,
}: SummonerHeroProps) {
  const [isComposerOpen, setIsComposerOpen] = useState(false);
  const [selectedLabel, setSelectedLabel] = useState<Label | null>(null);
  const previewLabels = playerLabels.slice(0, 3);
  const voteCountByLabelID = useMemo(
    () => new Map(playerLabels.map((label) => [label.id, label.voteCount])),
    [playerLabels],
  );

  function labelPercentage(voteCount: number) {
    if (totalLabelVotes === 0) {
      return 0;
    }
    return Math.round((voteCount / totalLabelVotes) * 100);
  }

  async function handleSaveLabel() {
    if (!selectedLabel || isSavingPlayerLabel) {
      return;
    }

    const saved = await onVotePlayerLabel(selectedLabel);
    if (!saved) {
      return;
    }

    setIsComposerOpen(false);
  }

  return (
    <section
      className={`rounded-2xl border p-6 shadow-xl shadow-black/20 transition ${
        isComposerOpen
          ? "border-cyan-300/35 bg-slate-900 ring-1 ring-cyan-300/20"
          : "border-white/10 bg-slate-900"
      }`}
    >
      <div className="flex flex-col gap-6">
        <div className="flex flex-col gap-5 md:flex-row md:items-center md:justify-between">
          <div className="flex flex-col gap-5 md:flex-row md:items-center">
            <div className="relative w-fit">
              {result.profileIconUrl ? (
                <Image
                  src={result.profileIconUrl}
                  alt={`${result.gameName} profile icon`}
                  width={112}
                  height={112}
                  className="rounded-xl border border-white/10"
                />
              ) : (
                <div className="flex h-28 w-28 items-center justify-center rounded-xl border border-dashed border-white/15 text-sm text-slate-400">
                  icon
                </div>
              )}

              <div className="absolute -bottom-2.5 left-1/2 -translate-x-1/2 rounded-lg border border-slate-500/40 bg-slate-700 px-2.5 py-0.5 text-lg font-semibold leading-none text-white shadow-lg shadow-black/25">
                {result.summonerLevel}
              </div>
            </div>

            <div className="min-w-0 space-y-3">
              <div className="flex flex-wrap items-center gap-2">
                {previewLabels.map((label) => (
                  <button
                    key={label.id}
                    type="button"
                    onClick={() => setIsComposerOpen(true)}
                    className="inline-flex items-center gap-2 rounded-full border border-[#5b3236] bg-[#201317] px-3 py-1.5 text-xs text-[#f5e8de] transition hover:border-[#7b474c] hover:bg-[#2a171c]"
                  >
                    <span>{label.name}</span>
                    <span className="text-[#f0b26f]">
                      {labelPercentage(label.voteCount)}%
                    </span>
                  </button>
                ))}
              </div>
              <h1 className="truncate pt-3 text-3xl font-semibold tracking-tight text-white md:pt-0 md:text-4xl">
                {result.gameName}#{result.tagLine}
              </h1>
              <p className="text-sm text-slate-400">{result.region}</p>
              {!isComposerOpen && labelError ? (
                <p className="text-sm text-red-300">{labelError}</p>
              ) : null}
            </div>
          </div>

          <button
            type="button"
            onClick={() => {
              setIsComposerOpen((current) => {
                if (current) {
                  setSelectedLabel(null);
                }

                return !current;
              });
            }}
            className={`inline-flex items-center justify-center rounded-xl px-5 py-3 text-sm font-semibold transition ${
              isComposerOpen
                ? "border border-cyan-300/30 bg-cyan-400/8 text-cyan-100 hover:bg-cyan-400/12"
                : "border border-cyan-300/30 bg-cyan-400/12 text-cyan-100 hover:border-cyan-300/45 hover:bg-cyan-400/18"
            }`}
          >
            {isComposerOpen ? "ラベル選択中" : "このプレイヤーにラベルを貼る"}
          </button>
        </div>

        {isComposerOpen ? (
          <div className="fade-up-soft border-t border-cyan-300/15 pt-6">
            <div className="rounded-2xl border border-cyan-300/20 bg-cyan-400/[0.07] p-5">
              <div className="flex flex-wrap gap-3">
                {availableLabels.map((label) => {
                  const isSelected = selectedLabel?.id === label.id;
                  const voteCount = voteCountByLabelID.get(label.id) ?? 0;

                  return (
                    <button
                      key={label.id}
                      type="button"
                      onClick={() => setSelectedLabel(label)}
                      className={`rounded-full px-4 py-2.5 text-sm font-medium transition ${
                        isSelected
                          ? "border border-[#f3c98b] bg-[#f0b26f] text-[#24150d] shadow-lg shadow-[#f0b26f]/20"
                          : "border border-[#5b3236] bg-[#201317] text-[#f0ddd1] hover:border-[#7b474c] hover:bg-[#2a171c] hover:text-white"
                      }`}
                    >
                      {label.name}
                      {voteCount > 0 ? (
                        <span className="ml-2 text-xs opacity-70">
                          {labelPercentage(voteCount)}%
                        </span>
                      ) : null}
                    </button>
                  );
                })}
              </div>

              {isLoadingPlayerLabels ? (
                <p className="mt-4 text-sm text-slate-300">ラベル集計を読み込んでいます。</p>
              ) : null}

              {selectedLabel ? (
                <div className="mt-5 flex flex-wrap items-center gap-3 border-t border-cyan-300/15 pt-4">
                  <span className="text-sm text-slate-300">選択中</span>
                  <span className="rounded-full bg-[#f0b26f] px-3 py-1 text-sm font-semibold text-[#24150d]">
                    {selectedLabel.name}
                  </span>
                  <button
                    type="button"
                    onClick={handleSaveLabel}
                    disabled={isSavingPlayerLabel}
                    className="rounded-lg border border-cyan-300/30 bg-cyan-400/12 px-4 py-2 text-sm font-semibold text-cyan-100 transition hover:bg-cyan-400/18 disabled:cursor-not-allowed disabled:opacity-60"
                  >
                    {isSavingPlayerLabel ? "保存中" : "保存"}
                  </button>
                </div>
              ) : null}
              {labelError ? <p className="mt-4 text-sm text-red-300">{labelError}</p> : null}
            </div>
          </div>
        ) : (
          <div className="fade-up-soft border-t border-white/8 pt-1" />
        )}
      </div>
    </section>
  );
}
