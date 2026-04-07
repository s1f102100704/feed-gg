import { SearchResult } from "@/lib/player-search";

import { RankCard } from "./rank-card";

type RankPanelProps = {
  result: SearchResult;
};

export function RankPanel({ result }: RankPanelProps) {
  return (
    <div className="space-y-4">
      <RankCard queueLabel="Solo Rank" rank={result.soloRank} />
      <RankCard queueLabel="Flex Rank" rank={result.flexRank} />
    </div>
  );
}
