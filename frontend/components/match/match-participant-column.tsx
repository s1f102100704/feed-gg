import Link from "next/link";
import Image from "next/image";

import {
  buildChampionIconUrl,
  formatParticipantName,
} from "@/lib/match-history";
import { buildSummonerPath } from "@/lib/player-search-path";
import { MatchParticipant, Region } from "@/types/player-search";

type MatchParticipantColumnProps = {
  assetVersion: string;
  participants: MatchParticipant[];
  region: Region;
  targetPUUID: string;
};

export function MatchParticipantColumn({
  assetVersion,
  participants,
  region,
  targetPUUID,
}: MatchParticipantColumnProps) {
  return (
    <div className="space-y-[1px]">
      {participants.map((participant) => {
        const isCurrentPlayer = participant.puuid === targetPUUID;
        const participantName = formatParticipantName(participant);
        const championIconUrl = buildChampionIconUrl(
          assetVersion,
          participant.championName,
        );
        const participantPath = buildSummonerPath(
          region,
          `${participant.gameName}#${participant.tagLine}`,
        );

        return (
          <div
            key={`${participant.puuid}-${participant.championName}`}
            className="flex items-center gap-1.5 rounded-md px-px py-px"
          >
            {championIconUrl ? (
              <Image
                src={championIconUrl}
                alt={participant.championName}
                width={22}
                height={22}
                className="h-[22px] w-[22px] rounded-[5px] border border-white/15 object-cover"
              />
            ) : (
              <div className="h-[22px] w-[22px] rounded-[5px] border border-white/15 bg-black/10" />
            )}

            {participantPath ? (
              <Link
                href={participantPath}
                className={`min-w-0 truncate text-xs leading-none hover:underline ${
                  isCurrentPlayer ? "font-semibold text-white" : "text-slate-100/85"
                }`}
                title={participantName}
              >
                {participantName}
              </Link>
            ) : (
              <p
                className={`min-w-0 truncate text-xs leading-none ${
                  isCurrentPlayer ? "font-semibold text-white" : "text-slate-100/85"
                }`}
                title={participantName}
              >
                {participantName}
              </p>
            )}
          </div>
        );
      })}
    </div>
  );
}
