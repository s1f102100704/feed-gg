"use client";

import { useMemo } from "react";
import { useParams } from "next/navigation";

import { SummonerScreen } from "@/features/summoner/summoner-screen";
import { pathSegmentToRegion } from "@/lib/player-search-path";

type SummonerPageParams = {
  region: string;
  gameName: string;
  tagLine: string;
};

export default function SummonerPage() {
  const params = useParams<SummonerPageParams>();
  const resolvedRegion = useMemo(
    () => pathSegmentToRegion(params.region),
    [params.region],
  );
  const decodedGameName = useMemo(
    () => decodeURIComponent(params.gameName),
    [params.gameName],
  );
  const decodedTagLine = useMemo(
    () => decodeURIComponent(params.tagLine),
    [params.tagLine],
  );
  const initialRiotId = `${decodedGameName}#${decodedTagLine}`;

  return (
    <SummonerScreen
      key={`${params.region}/${params.gameName}/${params.tagLine}`}
      resolvedRegion={resolvedRegion}
      decodedGameName={decodedGameName}
      decodedTagLine={decodedTagLine}
      initialRiotId={initialRiotId}
    />
  );
}
