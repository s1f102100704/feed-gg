import Image from "next/image";

import {
  buildChampionIconUrl,
  formatParticipantName,
} from "@/lib/match-history";
import { MatchParticipant } from "@/lib/player-search";

type MatchParticipantColumnProps = {
  assetVersion: string;
  participants: MatchParticipant[];
  targetPUUID: string;
};

export function MatchParticipantColumn({
  assetVersion,
  participants,
  targetPUUID,
}: MatchParticipantColumnProps) {
  return (
    <div className="space-y-[1px]">
      {participants.map((participant) => {
        const isCurrentPlayer = participant.puuid === targetPUUID;
        const championIconUrl = buildChampionIconUrl(
          assetVersion,
          participant.championName,
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

            <p
              className={`min-w-0 truncate text-xs leading-none ${
                isCurrentPlayer ? "font-semibold text-white" : "text-slate-100/85"
              }`}
              title={formatParticipantName(participant)}
            >
              {formatParticipantName(participant)}
            </p>
          </div>
        );
      })}
    </div>
  );
}
