import Image from "next/image";

type MatchChampionPortraitProps = {
  championName: string;
  imageUrl: string;
};

export function MatchChampionPortrait({
  championName,
  imageUrl,
}: MatchChampionPortraitProps) {
  if (!imageUrl) {
    return (
      <div className="flex h-[64px] w-[64px] items-center justify-center rounded-xl border border-white/15 bg-black/15 p-1.5 text-center text-[10px] text-slate-200/80">
        {championName}
      </div>
    );
  }

  return (
    <Image
      src={imageUrl}
      alt={championName}
      width={64}
      height={64}
      className="h-[64px] w-[64px] rounded-xl border border-white/15 object-cover shadow-lg shadow-black/25"
    />
  );
}
