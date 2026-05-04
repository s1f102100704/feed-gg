import { API_BASE_URL } from "@/lib/player-search-api";
import { ApiError } from "@/types/player-search";
import {
  PlayerLabelsResponse,
  PlayerLabelVoteResponse,
} from "@/types/player-labels";

export async function fetchPlayerLabels(
  puuid: string,
  signal?: AbortSignal,
): Promise<PlayerLabelsResponse | ApiError> {
  const response = await fetch(
    `${API_BASE_URL}/api/players/${encodeURIComponent(puuid)}/labels`,
    { signal },
  );
  const payload = (await response.json()) as PlayerLabelsResponse | ApiError;

  if (!response.ok) {
    return {
      error: "error" in payload ? payload.error : "プレイヤーのラベル取得に失敗しました。",
    };
  }

  return payload as PlayerLabelsResponse;
}

export async function votePlayerLabel(
  puuid: string,
  labelId: number,
): Promise<PlayerLabelVoteResponse | ApiError> {
  const response = await fetch(
    `${API_BASE_URL}/api/players/${encodeURIComponent(puuid)}/labels`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ labelId }),
    },
  );
  const payload = (await response.json()) as PlayerLabelVoteResponse | ApiError;

  if (!response.ok) {
    return {
      error: "error" in payload ? payload.error : "ラベルの保存に失敗しました。",
    };
  }

  return payload as PlayerLabelVoteResponse;
}
