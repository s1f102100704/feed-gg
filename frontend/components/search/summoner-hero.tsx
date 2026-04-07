import Image from "next/image";

import { SearchResult } from "@/lib/player-search";

type SummonerHeroProps = {
  result: SearchResult;
};

export function SummonerHero({ result }: SummonerHeroProps) {
  return (
    <section className="rounded-[32px] border border-white/10 bg-slate-900 p-6 shadow-xl shadow-black/20">
      <div className="flex flex-col gap-5 md:flex-row md:items-center">
        <div className="relative w-fit">
          {result.profileIconUrl ? (
            <Image
              src={result.profileIconUrl}
              alt={`${result.gameName} profile icon`}
              width={112}
              height={112}
              className="rounded-[28px] border border-white/10"
            />
          ) : (
            <div className="flex h-28 w-28 items-center justify-center rounded-[28px] border border-dashed border-white/15 text-sm text-slate-400">
              icon
            </div>
          )}

          <div className="absolute -bottom-2.5 left-1/2 -translate-x-1/2 rounded-lg border border-slate-500/40 bg-slate-700 px-2.5 py-0.5 text-lg font-semibold leading-none text-white shadow-lg shadow-black/25">
            {result.summonerLevel}
          </div>
        </div>

        <div className="min-w-0 space-y-3">
          <div className="flex flex-wrap items-center gap-3">
            <span className="rounded-full bg-cyan-400/12 px-3 py-1 text-xs font-medium text-cyan-200">
              {result.region}
            </span>
          </div>
          <h1 className="truncate pt-3 text-3xl font-semibold tracking-tight text-white md:pt-0 md:text-4xl">
            {result.gameName}#{result.tagLine}
          </h1>
        </div>
      </div>
    </section>
  );
}
