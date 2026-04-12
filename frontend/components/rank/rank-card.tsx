import { calculateWinRate, formatTierText } from "@/lib/rank";
import { RankedQueue } from "@/types/player-search";

import { RankEmblem } from "./rank-emblem";

type RankCardProps = {
  queueLabel: string;
  rank?: RankedQueue;
};

export function RankCard({ queueLabel, rank }: RankCardProps) {
  const winRate = calculateWinRate(rank);

  return (
    <section className="rounded-xl border border-white/10 bg-slate-900 p-5 shadow-xl shadow-black/15">
      <div className="mb-4">
        <p className="text-xs uppercase tracking-[0.28em] text-slate-400">{queueLabel}</p>
      </div>

      {rank ? (
        <div className="flex items-center justify-between gap-4">
          <div className="flex min-w-0 items-center gap-4">
            <div className="flex h-20 w-20 items-center justify-center rounded-full bg-black/15">
              <RankEmblem rank={rank} />
            </div>

            <div className="min-w-0">
              <p className="text-lg font-semibold leading-none tracking-tight text-white">
                {formatTierText(rank)}
              </p>
              <p className="mt-3 text-sm text-slate-400">{rank.leaguePoints} LP</p>
            </div>
          </div>

          <dl className="shrink-0 space-y-1.5 text-right">
            <div>
              <dd className="text-sm font-medium text-slate-400">
                {rank.wins}勝 {rank.losses}敗
              </dd>
            </div>
            <div>
              <dd className="text-sm font-medium text-slate-400">勝率 {winRate}%</dd>
            </div>
          </dl>
        </div>
      ) : (
        <div className="flex items-center gap-4">
          <div className="flex h-20 w-20 items-center justify-center rounded-full bg-black/15">
            <RankEmblem rank={rank} />
          </div>
          <p className="text-lg text-slate-400">まだランク情報がありません</p>
        </div>
      )}
    </section>
  );
}
