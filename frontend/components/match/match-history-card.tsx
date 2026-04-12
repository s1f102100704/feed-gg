import {
  buildChampionIconUrl,
  formatDuration,
  formatKDA,
  formatPlayedAt,
  formatQueueLabel,
  formatTimeAgo,
  groupParticipantsByTeam,
} from "@/lib/match-history";
import { MatchSummary } from "@/lib/player-search";

import { MatchChampionPortrait } from "./match-champion-portrait";
import { MatchParticipantColumn } from "./match-participant-column";

type MatchHistoryCardProps = {
  assetVersion: string;
  match: MatchSummary;
  targetPUUID: string;
};

export function MatchHistoryCard({
  assetVersion,
  match,
  targetPUUID,
}: MatchHistoryCardProps) {
  const { allies, enemies } = groupParticipantsByTeam(match, targetPUUID);

  return (
    <article
      className={`rounded-2xl border px-3 py-2.5 shadow-xl shadow-black/10 md:px-3.5 md:py-2.5 ${
        match.win
          ? "border-sky-300/25 bg-sky-900/55"
          : "border-rose-300/25 bg-rose-900/50"
      }`}
    >
      <div className="flex flex-col gap-3 lg:grid lg:grid-cols-[138px_200px_minmax(0,1fr)] lg:items-center lg:gap-4">
        <div className="space-y-1.5">
          <div>
            <p
              className={`text-lg font-semibold ${
                match.win ? "text-sky-300" : "text-rose-300"
              }`}
            >
              {formatQueueLabel(match)}
            </p>
            <p className="mt-0.5 text-sm text-slate-100">
              {formatTimeAgo(match.playedAt)}
            </p>
          </div>

          <div className="space-y-0.5">
            <p
              className={`text-base font-semibold ${
                match.win ? "text-sky-200" : "text-rose-200"
              }`}
            >
              {match.win ? "Win" : "Lose"}{" "}
              <span className="font-medium text-slate-100/85">
                {formatDuration(match.durationSeconds)}
              </span>
            </p>
            <p className="text-[11px] text-slate-200/70">
              {formatPlayedAt(match.playedAt)}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <MatchChampionPortrait
            championName={match.championName}
            imageUrl={buildChampionIconUrl(assetVersion, match.championName)}
          />

          <div className="min-w-0 space-y-0.5">
            <p className="text-[1.45rem] font-semibold leading-none text-white">
              {match.kills}
              <span className="px-0.5 text-rose-300">/</span>
              <span className="text-rose-200">{match.deaths}</span>
              <span className="px-0.5 text-rose-300">/</span>
              {match.assists}
            </p>
            <p className="text-sm text-slate-100/85">{formatKDA(match)}</p>
          </div>
        </div>

        <div className="grid gap-2 sm:grid-cols-2 lg:self-stretch">
          <MatchParticipantColumn
            assetVersion={assetVersion}
            participants={allies}
            targetPUUID={targetPUUID}
          />
          <MatchParticipantColumn
            assetVersion={assetVersion}
            participants={enemies}
            targetPUUID={targetPUUID}
          />
        </div>
      </div>
    </article>
  );
}
