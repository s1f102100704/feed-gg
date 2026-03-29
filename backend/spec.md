プレイヤー情報
まず account-v1 で Riot ID（gameName + tagLine）から puuid を取り、そのあと summoner-v4 を puuid で引く。
account-v1 の /riot/account/v1/accounts/by-riot-id/{gameName}/{tagLine} で puuid が取得でき、summoner-v4 の /lol/summoner/v4/summoners/by-puuid/{encryptedPUUID} でサモナー情報と summonerID が取得できる。Riotはプレイヤー指定には Riot ID を使い、使えるなら PUUID ベースのエンドポイントを推奨している。

対戦履歴情報
match-v5 を使う。
/lol/match/v5/matches/by-puuid/{puuid}/ids で試合ID一覧を取り、/lol/match/v5/matches/{matchId} で試合詳細を取る流れ。API一覧にも match-v5 はLoL向けAPIとして載っている。

プレイヤーアイコン情報
アイコン番号そのものは summoner-v4 のサモナー情報から profileIconId として取る。
実際の画像ファイルは Riot API 本体ではなく、LoLドキュメントにある Data Dragon の静的アセットを使う。LoLドキュメントでは Data Dragon がプロフィールアイコンを含むゲームデータと画像アセットを提供すると案内している。

プレイヤーのランク情報
league-v4 を使う。
プレイヤーのランクは summoner-v4 で取った summonerID を使って league-v4 を引く、という流れになる。API一覧に league-v4 があり、LoLドキュメントでもプレイヤーの絞り込みには PUUID または summonerID が使え、APIによってどちらを使うかが分かれると説明されている。