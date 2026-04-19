import "server-only";

import { Region } from "@/types/player-search";

const SERVER_API_BASE_URL =
  process.env.INTERNAL_API_BASE_URL ??
  process.env.NEXT_PUBLIC_API_BASE_URL ??
  "http://localhost:8080";

type RegionsResponse = {
  regions: Region[];
};

export async function fetchRegions(): Promise<Region[]> {
  const response = await fetch(`${SERVER_API_BASE_URL}/api/regions`, {
    next: { revalidate: 3600 },
  });

  if (!response.ok) {
    throw new Error("failed to load regions");
  }

  const payload = (await response.json()) as RegionsResponse;
  return payload.regions;
}
