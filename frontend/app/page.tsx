import { HomeScreen } from "@/features/home/home-screen";
import { SUPPORTED_REGIONS } from "@/lib/regions";

export default function Home() {
  return <HomeScreen initialRegions={SUPPORTED_REGIONS} />;
}
