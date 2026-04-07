export type Region =
  | "BR1"
  | "EUN1"
  | "EUW1"
  | "JP1"
  | "KR"
  | "LA1"
  | "LA2"
  | "ME1"
  | "NA1"
  | "OC1"
  | "PH2"
  | "RU"
  | "SG2"
  | "TH2"
  | "TR1"
  | "TW2"
  | "VN2";

export type SearchResult = {
  region: string;
  puuid: string;
  gameName: string;
  tagLine: string;
  summonerLevel: number;
  profileIconId: number;
  profileIconUrl: string;
  revisionDate: number;
  soloRank?: RankedQueue;
  flexRank?: RankedQueue;
};

export type RankedQueue = {
  tier: string;
  rank: string;
  leaguePoints: number;
  wins: number;
  losses: number;
};

export type ApiError = {
  error: string;
};

export const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export const regionOptions: Region[] = [
  "KR",
  "JP1",
  "NA1",
  "EUW1",
  "EUN1",
  "BR1",
  "LA1",
  "LA2",
  "ME1",
  "OC1",
  "PH2",
  "RU",
  "SG2",
  "TH2",
  "TR1",
  "TW2",
  "VN2",
];

export const rankEmblemImageMap = {
  IRON: "/Ranked%20Emblems%20Latest/Rank=Iron.png",
  BRONZE: "/Ranked%20Emblems%20Latest/Rank=Bronze.png",
  SILVER: "/Ranked%20Emblems%20Latest/Rank=Silver.png",
  GOLD: "/Ranked%20Emblems%20Latest/Rank=Gold.png",
  PLATINUM: "/Ranked%20Emblems%20Latest/Rank=Platinum.png",
  EMERALD: "/Ranked%20Emblems%20Latest/Rank=Emerald.png",
  DIAMOND: "/Ranked%20Emblems%20Latest/Rank=Diamond.png",
  MASTER: "/Ranked%20Emblems%20Latest/Rank=Master.png",
  GRANDMASTER: "/Ranked%20Emblems%20Latest/Rank=Grandmaster.png",
  CHALLENGER: "/Ranked%20Emblems%20Latest/Rank=Challenger.png",
} as const;

const regionPathSegmentMap: Record<Region, string> = {
  BR1: "br",
  EUN1: "eune",
  EUW1: "euw",
  JP1: "jp",
  KR: "kr",
  LA1: "lan",
  LA2: "las",
  ME1: "me",
  NA1: "na",
  OC1: "oce",
  PH2: "ph",
  RU: "ru",
  SG2: "sg",
  TH2: "th",
  TR1: "tr",
  TW2: "tw",
  VN2: "vn",
};

const pathSegmentRegionMap = Object.fromEntries(
  Object.entries(regionPathSegmentMap).map(([region, segment]) => [segment, region]),
) as Record<string, Region>;

export function parseRiotID(value: string) {
  const trimmed = value.trim();
  const separatorIndex = trimmed.lastIndexOf("#");

  if (separatorIndex <= 0 || separatorIndex === trimmed.length - 1) {
    return null;
  }

  const gameName = trimmed.slice(0, separatorIndex).trim();
  const tagLine = trimmed.slice(separatorIndex + 1).trim();

  if (!gameName || !tagLine) {
    return null;
  }

  return { gameName, tagLine };
}

export function regionToPathSegment(region: Region) {
  return regionPathSegmentMap[region];
}

export function pathSegmentToRegion(segment: string): Region | null {
  return pathSegmentRegionMap[segment.toLowerCase()] ?? null;
}

export function buildSummonerPath(region: Region, riotId: string) {
  const parsed = parseRiotID(riotId);
  if (!parsed) {
    return null;
  }

  return `/summoners/${regionToPathSegment(region)}/${encodeURIComponent(
    parsed.gameName,
  )}/${encodeURIComponent(parsed.tagLine)}`;
}

export function formatLastUpdated(revisionDate: number) {
  if (!revisionDate) {
    return "No activity timestamp";
  }

  return new Intl.DateTimeFormat("ja-JP", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(new Date(revisionDate));
}

export function formatTierText(rank?: RankedQueue) {
  if (!rank) {
    return "Unranked";
  }

  return `${rank.tier} ${rank.rank}`;
}

export function calculateWinRate(rank?: RankedQueue) {
  if (!rank) {
    return 0;
  }

  const totalGames = rank.wins + rank.losses;
  if (totalGames === 0) {
    return 0;
  }

  return Math.round((rank.wins / totalGames) * 100);
}

export function getRankEmblemImageSrc(rank?: RankedQueue) {
  if (!rank) {
    return null;
  }

  return rankEmblemImageMap[rank.tier as keyof typeof rankEmblemImageMap] ?? null;
}
