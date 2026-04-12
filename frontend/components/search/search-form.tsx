import { regionOptions } from "@/lib/player-search-path";
import { Region } from "@/types/player-search";

type SearchFormProps = {
  region: Region;
  riotId: string;
  isLoading: boolean;
  onRegionChange: (region: Region) => void;
  onRiotIDChange: (value: string) => void;
  onSubmit: (event: React.FormEvent<HTMLFormElement>) => void;
};

export function SearchForm({
  region,
  riotId,
  isLoading,
  onRegionChange,
  onRiotIDChange,
  onSubmit,
}: SearchFormProps) {
  return (
    <form
      onSubmit={onSubmit}
      className="mx-auto flex w-full max-w-3xl flex-col gap-3 rounded-xl border border-white/10 bg-white/5 p-3 shadow-2xl shadow-black/20 backdrop-blur md:flex-row"
    >
      <select
        value={region}
        onChange={(event) => onRegionChange(event.target.value as Region)}
        className="h-14 rounded-xl border border-white/10 bg-slate-950/80 px-4 text-sm text-white outline-none md:w-32"
      >
        {regionOptions.map((option) => (
          <option key={option} value={option} className="bg-slate-950">
            {option}
          </option>
        ))}
      </select>

      <input
        value={riotId}
        onChange={(event) => onRiotIDChange(event.target.value)}
        placeholder="プレイヤー名#JP1"
        className="h-14 min-w-0 flex-1 rounded-xl border border-white/10 bg-slate-950/80 px-5 text-base text-white placeholder:text-slate-500 outline-none"
      />

      <button
        type="submit"
        disabled={isLoading}
        className="h-14 rounded-xl bg-cyan-400 px-6 text-sm font-semibold text-slate-950 transition hover:bg-cyan-300 disabled:cursor-not-allowed disabled:opacity-60 md:w-36"
      >
        {isLoading ? "Searching..." : "Search"}
      </button>
    </form>
  );
}
