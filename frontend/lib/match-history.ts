import {
  MatchParticipant,
  MatchSummary,
  SearchResult,
} from "@/types/player-search";

const queueLabelMap: Record<number, string> = {
  400: "ノーマルドラフト",
  420: "ソロランク",
  430: "ノーマルブラインド",
  440: "フレックス",
  450: "ARAM",
  490: "クイックプレイ",
  700: "Clash",
  830: "AI Intro",
  840: "AI Beginner",
  850: "AI Intermediate",
  900: "URF",
  1700: "アリーナ",
};

const championIdOverrides: Record<string, string> = {
  FiddleSticks: "Fiddlesticks",
  Wukong: "MonkeyKing",
};

const roleOrder: Record<string, number> = {
  TOP: 0,
  JUNGLE: 1,
  MIDDLE: 2,
  BOTTOM: 3,
  UTILITY: 4,
  SUPPORT: 4,
  UNKNOWN: 5,
};

export function formatQueueLabel(match: MatchSummary) {
  return queueLabelMap[match.queueId] ?? formatGameMode(match.gameMode);
}

export function formatGameMode(gameMode: string) {
  switch (gameMode) {
    case "CLASSIC":
      return "サモナーズリフト";
    case "ARAM":
      return "ARAM";
    case "URF":
      return "URF";
    case "CHERRY":
      return "アリーナ";
    default:
      return gameMode.replace(/_/g, " ");
  }
}

export function formatTimeAgo(timestamp: number, now = Date.now()) {
  const diffMs = Math.max(now - timestamp, 0);
  const minute = 60 * 1000;
  const hour = 60 * minute;
  const day = 24 * hour;

  if (diffMs < minute) {
    return "たった今";
  }
  if (diffMs < hour) {
    return `${Math.floor(diffMs / minute)}分前`;
  }
  if (diffMs < day) {
    return `${Math.floor(diffMs / hour)}時間前`;
  }
  if (diffMs < 7 * day) {
    return `${Math.floor(diffMs / day)}日前`;
  }

  return new Intl.DateTimeFormat("ja-JP", {
    month: "numeric",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(timestamp));
}

export function formatPlayedAt(timestamp: number) {
  return new Intl.DateTimeFormat("ja-JP", {
    year: "numeric",
    month: "numeric",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(timestamp));
}

export function formatDuration(durationSeconds: number) {
  const totalSeconds = Math.max(Math.floor(durationSeconds), 0);
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;

  if (hours > 0) {
    return `${hours}:${String(minutes).padStart(2, "0")}:${String(seconds).padStart(
      2,
      "0",
    )}`;
  }

  return `${minutes}:${String(seconds).padStart(2, "0")}`;
}

export function formatKDA(match: MatchSummary) {
  if (match.deaths === 0) {
    return "Perfect KDA";
  }

  return `${((match.kills + match.assists) / match.deaths).toFixed(2)} KDA`;
}

export function buildChampionIconUrl(version: string, championName: string) {
  const championId = championIdOverrides[championName] ?? championName;
  if (!version || !championId) {
    return "";
  }

  return `https://ddragon.leagueoflegends.com/cdn/${version}/img/champion/${encodeURIComponent(
    championId,
  )}.png`;
}

export function resolveDataDragonVersion(result: SearchResult) {
  const versionFromProfile = result.profileIconUrl.match(/\/cdn\/([^/]+)\//)?.[1];
  if (versionFromProfile) {
    return versionFromProfile;
  }

  return normalizeGameVersion(result.matches[0]?.gameVersion ?? "") ?? "";
}

export function groupParticipantsByTeam(match: MatchSummary, targetPUUID: string) {
  const currentParticipant =
    match.participants.find((participant) => participant.puuid === targetPUUID) ?? null;
  const allyWinState = currentParticipant?.win ?? match.win;

  return {
    allies: sortParticipants(
      match.participants.filter((participant) => participant.win === allyWinState),
    ),
    enemies: sortParticipants(
      match.participants.filter((participant) => participant.win !== allyWinState),
    ),
  };
}

export function formatParticipantName(participant: MatchParticipant) {
  if (participant.gameName.trim()) {
    return participant.gameName;
  }
  if (participant.tagLine.trim()) {
    return `#${participant.tagLine}`;
  }
  return "Unknown";
}

function normalizeGameVersion(gameVersion: string) {
  const [major, minor] = gameVersion.split(".");
  if (!major || !minor) {
    return null;
  }

  return `${major}.${minor}.1`;
}

function sortParticipants(participants: MatchParticipant[]) {
  return [...participants].sort((left, right) => {
    const roleDiff = (roleOrder[left.role] ?? 99) - (roleOrder[right.role] ?? 99);
    if (roleDiff !== 0) {
      return roleDiff;
    }

    return formatParticipantName(left).localeCompare(
      formatParticipantName(right),
      "ja",
    );
  });
}
