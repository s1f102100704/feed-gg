import { RankedQueue } from "@/types/player-search";

const rankEmblemImageMap = {
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
