"use client";

import Image from "next/image";
import { useState } from "react";

import { Label, SearchResult } from "@/types/player-search";

type SummonerHeroProps = {
  labels: Label[];
  result: SearchResult;
};

export function SummonerHero({ result, labels }: SummonerHeroProps) {
  const [isComposerOpen, setIsComposerOpen] = useState(false);
  const [selectedLabel, setSelectedLabel] = useState<Label | null>(null);
  const previewLabels = labels.slice(0, 3);

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
                  </button>
                ))}
              </div>
              <h1 className="truncate pt-3 text-3xl font-semibold tracking-tight text-white md:pt-0 md:text-4xl">
                {result.gameName}#{result.tagLine}
              </h1>
              <p className="text-sm text-slate-400">{result.region}</p>
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
                {labels.map((label) => {
                  const isSelected = selectedLabel?.id === label.id;

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
                    </button>
                  );
                })}
              </div>

              {selectedLabel ? (
                <div className="mt-5 flex flex-wrap items-center gap-3 border-t border-cyan-300/15 pt-4">
                  <span className="text-sm text-slate-300">選択中</span>
                  <span className="rounded-full bg-[#f0b26f] px-3 py-1 text-sm font-semibold text-[#24150d]">
                    {selectedLabel.name}
                  </span>
                </div>
              ) : null}
            </div>
          </div>
        ) : (
          <div className="fade-up-soft border-t border-white/8 pt-1" />
        )}
      </div>
    </section>
  );
}
