import Image from "next/image";

import { formatLastUpdated, SearchResult } from "@/lib/player-search";

type SearchResultCardProps = {
  result: SearchResult;
};

export function SearchResultCard({ result }: SearchResultCardProps) {
  return (
    <section className="mx-auto flex max-w-xl items-center gap-4 rounded-[28px] border border-white/10 bg-white/5 p-5 text-left backdrop-blur">
      {result.profileIconUrl ? (
        <Image
          src={result.profileIconUrl}
          alt={`${result.gameName} profile icon`}
          width={80}
          height={80}
          className="rounded-3xl border border-white/10"
        />
      ) : (
        <div className="flex h-20 w-20 items-center justify-center rounded-3xl border border-dashed border-white/15 text-sm text-slate-400">
          icon
        </div>
      )}

      <div className="min-w-0 space-y-1">
        <p className="truncate text-2xl font-semibold">
          {result.gameName}#{result.tagLine}
        </p>
        <p className="text-sm text-slate-400">
          {result.region} · Level {result.summonerLevel}
        </p>
        <p className="text-sm text-slate-500">
          Updated {formatLastUpdated(result.revisionDate)}
        </p>
      </div>
    </section>
  );
}
