import { PATH_SEGMENT_REGIONS, REGION_PATH_SEGMENTS, Region } from "@/lib/regions";

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
  return REGION_PATH_SEGMENTS[region];
}

export function pathSegmentToRegion(segment: string): Region | null {
  return PATH_SEGMENT_REGIONS[segment.toLowerCase()] ?? null;
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

export function buildPlayerSearchApiPath(
  region: Region,
  gameName: string,
  tagLine: string,
) {
  return `/api/players/${encodeURIComponent(region)}/${encodeURIComponent(
    gameName,
  )}/${encodeURIComponent(tagLine)}`;
}
