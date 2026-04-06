import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "feed.gg",
  description: "League of Legends player search prototype",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja">
      <body>{children}</body>
    </html>
  );
}
