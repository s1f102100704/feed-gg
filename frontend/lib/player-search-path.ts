import { Region } from "@/types/player-search";

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
