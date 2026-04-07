import { ReactNode } from "react";

type SummonerContentLayoutProps = {
  left: ReactNode;
  right: ReactNode;
};

export function SummonerContentLayout({
  left,
  right,
}: SummonerContentLayoutProps) {
  return (
    <section className="grid gap-6 lg:grid-cols-[340px_minmax(0,1fr)]">
      <aside>{left}</aside>
      <section>{right}</section>
    </section>
  );
}
