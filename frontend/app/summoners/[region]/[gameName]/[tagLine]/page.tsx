import { SummonerScreen } from "@/features/summoner/summoner-screen";
import { pathSegmentToRegion } from "@/lib/player-search-path";
import { SUPPORTED_REGIONS } from "@/lib/regions";

type SummonerPageProps = {
  params: Promise<{
    region: string;
    gameName: string;
    tagLine: string;
  }>;
};

export default async function SummonerPage({ params }: SummonerPageProps) {
  const { region, gameName, tagLine } = await params;
  const resolvedRegion = pathSegmentToRegion(region);
  const decodedGameName = decodeURIComponent(gameName);
  const decodedTagLine = decodeURIComponent(tagLine);
  const initialRiotId = `${decodedGameName}#${decodedTagLine}`;

  return (
    <SummonerScreen
      key={`${region}/${gameName}/${tagLine}`}
      initialRegions={SUPPORTED_REGIONS}
      resolvedRegion={resolvedRegion}
      decodedGameName={decodedGameName}
      decodedTagLine={decodedTagLine}
      initialRiotId={initialRiotId}
    />
  );
}
