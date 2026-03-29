# Riot API 仕様メモ

更新日: 2026-03-24

このファイルは `ReadMe.md` に書かれている要件を実装するために、今回のアプリで必要な Riot API と関連静的データを公式ドキュメント基準で整理したものです。

## 1. 今回のアプリで必要な API

必須:

1. `account-v1`
   - Riot ID (`gameName` + `tagLine`) から `puuid` を引く
2. `summoner-v4`
   - `puuid` から現在のサモナー概要を引く
3. `league-v4`
   - `puuid` からランク情報を引く
4. `match-v5`
   - `puuid` から試合 ID 一覧を引き、各試合詳細を引く
5. Data Dragon
   - プロフィールアイコン画像を引く

任意:

1. `account-v1` の active region API
   - 入力された LoL リージョンと実際の所属リージョンの突き合わせに使える
2. `queues.json`
   - `queueId` を UI 表示用ラベルに変換するために使える

## 2. 推奨の取得フロー

### 2.1 プレイヤー検索

1. フロントから `region`, `gameName`, `tagLine` を受け取る
2. `account-v1` で Riot ID から `puuid` を取得する
3. `summoner-v4` で `puuid` からプロフィール情報を取得する
4. `league-v4` で `puuid` からランク情報を取得する
5. `match-v5` で `puuid` から match id 一覧を取得する
6. `match-v5` で各 match id の試合詳細を取得する
7. `summoner-v4.profileIconId` と Data Dragon を使ってアイコン URL を組み立てる
8. 取得した内容を DB に保存する

### 2.2 更新ボタン

プロフィール画面の更新ボタンは、上と同じ順序で Riot API を再実行し、DB を上書き更新すればよいです。

## 3. API ごとの仕様

### 3.1 account-v1

用途:
- Riot ID から `puuid` を得る
- 必要なら `puuid` から Riot ID を再取得する
- 必要なら LoL の active region を得る

ルーティング:
- regional routing
- `AMERICAS`, `ASIA`, `EUROPE`
- 公式には「どの cluster からでも account を引けるので、最寄り cluster を使う」方針

#### 必須エンドポイント

`GET /riot/account/v1/accounts/by-riot-id/{gameName}/{tagLine}`

主な path param:
- `gameName`
- `tagLine`

レスポンス `AccountDto`:
- `puuid: string`
- `gameName: string`
- `tagLine: string`

実装メモ:
- Riot ID 検索の入口はこの API にする
- `gameName` または `tagLine` は account によってレスポンスから省略される場合がある

#### 補助エンドポイント

`GET /riot/account/v1/accounts/by-puuid/{puuid}`

用途:
- DB に保存済みの `puuid` から Riot ID 表示名を補完・更新する

#### 任意エンドポイント

`GET /riot/account/v1/region/by-game/{game}/by-puuid/{puuid}`

用途:
- `game=lol` で active region を取得する
- 入力された region が正しいかの検証に使える

レスポンス `AccountRegionDTO`:
- `puuid`
- `game`
- `region`

### 3.2 summoner-v4

用途:
- 現在のプロフィールアイコン ID
- 最終更新時刻
- サモナーレベル

ルーティング:
- platform routing
- 例: `JP1`, `NA1`, `KR`, `EUW1`

#### 必須エンドポイント

`GET /lol/summoner/v4/summoners/by-puuid/{encryptedPUUID}`

主な path param:
- `encryptedPUUID`

レスポンス `SummonerDTO`:
- `profileIconId: int`
- `revisionDate: long`
- `puuid: string`
- `summonerLevel: long`

実装メモ:
- この DTO だけでプロフィール表示に必要な現在アイコンとレベルが取れる
- 今回は `league-v4` 側に `by-puuid` があるため、ランク取得のために別の ID へ変換する必要はない

### 3.3 league-v4

用途:
- ランク情報
- Solo/Duo と Flex を含む各 queue の現在ランク

ルーティング:
- platform routing
- 例: `JP1`, `NA1`, `KR`, `EUW1`

#### 必須エンドポイント

`GET /lol/league/v4/entries/by-puuid/{encryptedPUUID}`

主な path param:
- `encryptedPUUID`

レスポンス `Set[LeagueEntryDTO]`:
- `puuid: string`
- `queueType: string`
- `tier: string`
- `rank: string`
- `leaguePoints: int`
- `wins: int`
- `losses: int`
- `hotStreak: boolean`
- `veteran: boolean`
- `freshBlood: boolean`
- `inactive: boolean`
- `miniSeries`

実装メモ:
- 1 プレイヤーに対して複数 queue が返る
- UI ではまず `RANKED_SOLO_5x5` と `RANKED_FLEX_SR` を優先表示すると扱いやすい
- 昇格戦中は `miniSeries` が入ることがある

### 3.4 match-v5

用途:
- 対戦履歴の一覧取得
- 試合詳細取得

ルーティング:
- regional routing
- `AMERICAS`, `ASIA`, `EUROPE`, `SEA`

公式メモ:
- `AMERICAS` は NA, BR, LAN, LAS
- `ASIA` は KR, JP
- `EUROPE` は EUNE, EUW, ME1, TR, RU
- `SEA` は OCE, SG2, TW2, VN2

#### 必須エンドポイント 1

`GET /lol/match/v5/matches/by-puuid/{puuid}/ids`

主な path param:
- `puuid`

主な query param:
- `startTime: long` 秒
- `endTime: long` 秒
- `queue: int`
- `type: string`
- `start: int` 既定値 `0`
- `count: int` 既定値 `20`, 有効範囲 `0..100`

レスポンス:
- `List[string]` match id 一覧

実装メモ:
- 初回は `count=20` から始めると扱いやすい
- ランク戦だけ欲しいなら `type=ranked` か `queue` 指定で絞る

#### 必須エンドポイント 2

`GET /lol/match/v5/matches/{matchId}`

主な path param:
- `matchId`

レスポンス `MatchDto`:
- `metadata`
- `info`

`MetadataDto` の主要項目:
- `dataVersion`
- `matchId`
- `participants` (`puuid` 一覧)

`InfoDto` の主要項目:
- `gameCreation`
- `gameStartTimestamp`
- `gameEndTimestamp`
- `gameDuration`
- `gameVersion`
- `gameMode`
- `gameType`
- `mapId`
- `queueId`
- `platformId`
- `participants`
- `teams`
- `tournamentCode`

`ParticipantDto` のうち今回保存価値が高い項目:
- `puuid`
- `riotIdGameName`
- `riotIdTagline`
- `profileIcon`
- `championName`
- `teamPosition`
- `kills`
- `deaths`
- `assists`
- `win`
- `summoner1Id`
- `summoner2Id`

実装メモ:
- `gameVersion` の先頭 2 要素でパッチ判定しやすい
- `queueId` は UI 表示時に `queues.json` で名前変換するとよい
- `teamPosition` は Riot が `individualPosition` より推奨している
- `gameDuration` は時期によって単位の扱いが異なるため、`gameEndTimestamp` がある場合は秒、ない場合はミリ秒として扱う前提で正規化したほうが安全

### 3.5 Data Dragon

用途:
- 現在プロフィールのアイコン画像 URL を組み立てる

#### 必須エンドポイント

`GET https://ddragon.leagueoflegends.com/api/versions.json`

用途:
- 利用可能な Data Dragon バージョン一覧を取得する

#### 推奨補助エンドポイント

`GET https://ddragon.leagueoflegends.com/realms/{region}.json`

用途:
- 地域ごとの実クライアントバージョンを確認する
- JP なら `realms/jp.json`

#### 画像 URL パターン

`https://ddragon.leagueoflegends.com/cdn/{version}/img/profileicon/{profileIconId}.png`

実装メモ:
- `profileIconId` は `summoner-v4` から取る
- バージョンは `versions.json` の先頭を使うか、より厳密にやるなら `realms/{region}.json` を使う

## 4. このアプリで保存したい最小データ

### 4.1 players

- `puuid`
- `platform_region` 例: `JP1`
- `game_name`
- `tag_line`
- `profile_icon_id`
- `summoner_level`
- `revision_date`
- `last_synced_at`

### 4.2 ranked_entries

- `puuid`
- `queue_type`
- `tier`
- `rank`
- `league_points`
- `wins`
- `losses`
- `hot_streak`
- `veteran`
- `fresh_blood`
- `inactive`
- `mini_series_progress`
- `mini_series_target`
- `mini_series_wins`
- `mini_series_losses`

### 4.3 matches

- `match_id`
- `data_version`
- `platform_id`
- `queue_id`
- `map_id`
- `game_creation`
- `game_start_timestamp`
- `game_end_timestamp`
- `game_duration_seconds`
- `game_mode`
- `game_type`
- `game_version`
- `tournament_code`

### 4.4 match_participants

- `match_id`
- `puuid`
- `riot_id_game_name`
- `riot_id_tagline`
- `profile_icon`
- `champion_name`
- `team_position`
- `kills`
- `deaths`
- `assists`
- `win`
- `summoner1_id`
- `summoner2_id`

## 5. リージョンの扱い

今回の実装では、入力された LoL リージョンを 2 種類に分けて扱う必要があります。

1. platform routing 用
   - `summoner-v4`
   - `league-v4`
   - 例: `JP1`
2. regional routing 用
   - `account-v1`
   - `match-v5`
   - 例: `ASIA`

JP1 の場合の実装例:

- `summoner-v4`, `league-v4`: `https://jp1.api.riotgames.com`
- `match-v5`: `https://asia.api.riotgames.com`
- `account-v1`: `https://asia.api.riotgames.com` で問題なし

## 6. エラー処理とレート制限

今回使う API は共通して次のエラーを返しうる想定です。

- `400` Bad request
- `401` Unauthorized
- `403` Forbidden
- `404` Data not found
- `429` Rate limit exceeded
- `500`, `502`, `503`, `504`

Riot の portal ドキュメント上の personal key 既定レート制限:

- 1 秒あたり 20 リクエスト
- 2 分あたり 100 リクエスト
- いずれも region ごとに適用

実装メモ:
- `429` を受けたら `Retry-After` を見て待つ
- match detail をまとめて取りに行く時は並列数を制限する
- 「既に DB にいるプレイヤーをそのまま表示する」仕様でも、更新ボタン時は都度再同期する

## 7. 実装上の結論

今回のアプリに対しては、最低限この組み合わせで足ります。

1. Riot ID 検索: `account-v1`
2. プレイヤー概要: `summoner-v4`
3. ランク: `league-v4` の `by-puuid`
4. 対戦履歴一覧: `match-v5` の `by-puuid/{puuid}/ids`
5. 対戦詳細: `match-v5` の `matches/{matchId}`
6. アイコン画像: Data Dragon

つまり、README の要件を満たすために本当に必要な外部取得元は `account-v1`, `summoner-v4`, `league-v4`, `match-v5`, Data Dragon の 5 つです。

## 8. 参照元

- [Riot API Reference](https://developer.riotgames.com/apis)
- [League of Legends Docs](https://developer.riotgames.com/docs/lol)
- [Developer Portal Docs](https://developer.riotgames.com/docs/portal)
- [account-v1 detail](https://developer.riotgames.com/apis#account-v1)
- [summoner-v4 detail](https://developer.riotgames.com/apis#summoner-v4)
- [league-v4 detail](https://developer.riotgames.com/apis#league-v4)
- [match-v5 detail](https://developer.riotgames.com/apis#match-v5)
- [Data Dragon versions](https://ddragon.leagueoflegends.com/api/versions.json)
