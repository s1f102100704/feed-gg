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
  matches: MatchSummary[];
};

export type RankedQueue = {
  tier: string;
  rank: string;
  leaguePoints: number;
  wins: number;
  losses: number;
};

export type MatchParticipant = {
  puuid: string;
  gameName: string;
  tagLine: string;
  championName: string;
  role: string;
  win: boolean;
  kills: number;
  deaths: number;
  assists: number;
  summonerSpell1Id: number;
  summonerSpell2Id: number;
};

export type MatchSummary = {
  matchId: string;
  playedAt: number;
  gameVersion: string;
  gameMode: string;
  queueId: number;
  championName: string;
  role: string;
  win: boolean;
  kills: number;
  deaths: number;
  assists: number;
  summonerSpell1Id: number;
  summonerSpell2Id: number;
  durationSeconds: number;
  participants: MatchParticipant[];
};

export type ApiError = {
  error: string;
};
