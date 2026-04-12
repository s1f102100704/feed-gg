import { resolveDataDragonVersion } from "@/lib/match-history";
import { SearchResult } from "@/lib/player-search";

import { MatchHistoryCard } from "./match-history-card";

type MatchHistoryPanelProps = {
  result: SearchResult;
};

export function MatchHistoryPanel({ result }: MatchHistoryPanelProps) {
  if (!result.matches.length) {
    return (
      <section className="rounded-2xl border border-white/10 bg-slate-900 p-6 shadow-xl shadow-black/15">
        <div className="space-y-1">
          <p className="text-xs uppercase tracking-[0.28em] text-slate-400">
            Match History
          </p>
          <h2 className="text-2xl font-semibold tracking-tight text-white">
            対戦履歴
          </h2>
        </div>
        <p className="mt-4 text-sm text-slate-400">
          取得できる試合履歴がまだありません。
        </p>
      </section>
    );
  }

  const assetVersion = resolveDataDragonVersion(result);

  return (
    <section className="space-y-4">
      <div className="rounded-2xl border border-white/10 bg-slate-900 p-6 shadow-xl shadow-black/15">
        <div className="flex flex-wrap items-end justify-between gap-3">
          <div className="space-y-1">
            <p className="text-xs uppercase tracking-[0.28em] text-slate-400">
              Match History
            </p>
            <h2 className="text-2xl font-semibold tracking-tight text-white">
              対戦履歴
            </h2>
          </div>
          <p className="text-sm text-slate-400">{result.matches.length} games</p>
        </div>
      </div>

      <div className="space-y-3">
        {result.matches.map((match) => (
          <MatchHistoryCard
            key={match.matchId}
            assetVersion={assetVersion}
            match={match}
            targetPUUID={result.puuid}
          />
        ))}
      </div>
    </section>
  );
}
