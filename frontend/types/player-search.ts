import type { Region } from "@/lib/regions";

export type { Region };

export type SearchResult = {
  region: Region;
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

export type Label = {
  id: number;
  name: string;
};

export type PlayerLabelSummary = Label & {
  voteCount: number;
};

export type PlayerLabelsResponse = {
  labels: PlayerLabelSummary[];
  totalVotes: number;
};

export type PlayerLabelVoteResponse = PlayerLabelsResponse & {
  selectedLabel: PlayerLabelSummary;
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
  teamId: number;
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
  teamId: number;
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
