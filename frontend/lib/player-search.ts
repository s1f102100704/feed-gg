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

export function formatLastUpdated(revisionDate: number) {
  if (!revisionDate) {
    return "No activity timestamp";
  }

  return new Intl.DateTimeFormat("ja-JP", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(new Date(revisionDate));
}
