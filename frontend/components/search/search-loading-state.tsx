type SearchLoadingStateProps = {
  message: string;
};

export function SearchLoadingState({ message }: SearchLoadingStateProps) {
  return (
    <section className="rounded-[32px] border border-white/10 bg-white/5 p-8 text-slate-300">
      {message}
    </section>
  );
}
