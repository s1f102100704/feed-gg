import Image from "next/image";

import { getRankEmblemImageSrc, RankedQueue } from "@/lib/player-search";

type RankEmblemProps = {
  rank?: RankedQueue;
};

export function RankEmblem({ rank }: RankEmblemProps) {
  const imageSrc = getRankEmblemImageSrc(rank);

  return (
    <div className="flex h-24 w-24 items-center justify-center">
      {imageSrc ? (
        <Image
          src={imageSrc}
          alt={`${rank?.tier ?? "Unranked"} emblem`}
          width={96}
          height={96}
          className="h-24 w-24 object-contain"
        />
      ) : (
        <div className="flex h-20 w-20 items-center justify-center rounded-full border border-white/10 bg-slate-800 text-sm font-semibold text-slate-300">
          -
        </div>
      )}
    </div>
  );
}
