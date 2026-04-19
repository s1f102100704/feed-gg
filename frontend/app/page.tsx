import { HomeScreen } from "@/features/home/home-screen";
import { fetchRegions } from "@/lib/server-regions";

export default async function Home() {
  const regions = await fetchRegions();

  return <HomeScreen initialRegions={regions} />;
}
