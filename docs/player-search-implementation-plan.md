# Player Search Implementation Plan

## Goal

プレイヤー検索の最終フローは以下を目指す。

1. `cache` を確認する
2. `DB` を確認する
3. `Riot API` を叩く
4. 取得結果を `DB` に保存する
5. 検索結果を `cache` に保存する

検索キーは `region + gameName + tagLine` を使う。  
usecase の返り値は当面 `PlayerSearchResult` を正規形として扱う。

## Current Design

- `main.go`
  起動・依存の組み立て・routing を担当する
- `adapter/http/player_search.go`
  HTTP request/response を担当する
- `usecase/player_search.go`
  検索フローの分岐を担当する
- `infrastructure/riot`
  Riot API 取得を担当する
- `infrastructure/playersearch`
  DB 保存・取得を担当する
- `infrastructure/cache`
  キャッシュ実装を担当する

## Completed Prerequisites

以下は検索フロー本実装の前提として対応済み。

- `region` と `player_ranks` の seed migration を追加済み
- `player(region_id, game_name, tag_line)` の unique 制約へ変更済み
- `sqlc` query を責務単位で分割済み
  - `regions.sql`
  - `player_ranks.sql`
  - `players.sql`
  - `player_current_ranks.sql`
  - `matches.sql`
  - `match_participants.sql`
- `match_participant.team_id` を保存できるよう、Riot DTO / usecase / frontend の型へ `teamId` を追加済み
- SQL 側の trim / upper / NULLIF 正規化は削除済み
  - Go 側へ移す TODO は [TODO.md](/Users/ellery/dev/koshinankin/feed-gg/docs/TODO.md) に記載済み

## Step Status

### Step 1: Search Flow Interfaces

状態: 完了

目的:

- `usecase` の依存を `Cache / Repository / RiotGateway / RegionChecker` に整理する
- 後続ステップで `cache -> db -> riot -> save -> cache` を実装できる差し込み口を作る

実装済み内容:

- [backend/internal/usecase/player_search.go](/Users/ellery/dev/koshinankin/feed-gg/backend/internal/usecase/player_search.go)
  - `PlayerSearchCache` を追加
  - `PlayerSearchRepository` を追加
  - `RiotGateway` を追加
  - `PlayerSearchKey` と `NewPlayerSearchKey` を追加
  - `PlayerSearch` struct に `cache`, `repository`, `riotGateway`, `regionChecker` を保持する形へ変更
- [backend/internal/infrastructure/cache/player_search.go](/Users/ellery/dev/koshinankin/feed-gg/backend/internal/infrastructure/cache/player_search.go)
  - `NoopPlayerSearchCache` を追加
- [backend/internal/infrastructure/playersearch/repository.go](/Users/ellery/dev/koshinankin/feed-gg/backend/internal/infrastructure/playersearch/repository.go)
  - `NoopRepository` を追加
- [backend/cmd/api/main.go](/Users/ellery/dev/koshinankin/feed-gg/backend/cmd/api/main.go)
  - no-op の `cache` / `repository` を usecase に配線

この時点で意図的に未実装のもの:

- `Execute` 内での `cache.Get/Set`
- `Execute` 内での `repository.FindSavedPlayer`
- `Execute` 内での `repository.SaveFetchedPlayer`

現状の `Execute` はまだ以下の流れのまま。

1. `region` マスタ確認
2. `Riot API` 取得
3. `PlayerSearchResult` に変換して返す

### Step 2: DB Save

状態: 未着手

目的:

- `Riot API` から取得した 1 プレイヤー分を 1 transaction で保存できるようにする

実装要件:

- `PlayerSearchRepository.SaveFetchedPlayer(ctx, fetched)` を本実装にする
- transaction は `Riot API` 取得後に開始する
- `sqlc` の `Queries.WithTx(tx)` を使う
- 以下を順に保存する
  1. `region` を解決する
  2. `player` を upsert する
  3. `player_current_rank` を入れ直す
  4. `match_history` を upsert する
  5. `match_participant` 用に参加者の `player` を最小情報で upsert する
  6. `match_participant` を upsert する

使う query:

- [backend/query/regions.sql](/Users/ellery/dev/koshinankin/feed-gg/backend/query/regions.sql)
- [backend/query/player_ranks.sql](/Users/ellery/dev/koshinankin/feed-gg/backend/query/player_ranks.sql)
- [backend/query/players.sql](/Users/ellery/dev/koshinankin/feed-gg/backend/query/players.sql)
- [backend/query/player_current_ranks.sql](/Users/ellery/dev/koshinankin/feed-gg/backend/query/player_current_ranks.sql)
- [backend/query/matches.sql](/Users/ellery/dev/koshinankin/feed-gg/backend/query/matches.sql)
- [backend/query/match_participants.sql](/Users/ellery/dev/koshinankin/feed-gg/backend/query/match_participants.sql)

完了条件:

- `SaveFetchedPlayer(ctx, fetched)` で DB 保存が最後まで通る

### Step 3: DB Read

状態: 未着手

目的:

- `region + gameName + tagLine` から保存済みプレイヤーを組み立てて返せるようにする

実装要件:

- `PlayerSearchRepository.FindSavedPlayer(ctx, input)` を本実装にする
- 入口は `GetSavedPlayerKeyByRiotID` で `player_id`, `puuid` を引く
- その後、詳細を以下から集約する
  - `GetSavedPlayerByPuuid`
  - `ListPlayerCurrentRanksByPlayerID`
  - `ListRecentMatchHistoriesByPlayerID`
  - `ListMatchParticipantsByMatchHistoryID`
- DB の row shape ではなく `PlayerSearchResult` を返す

完了条件:

- `FindSavedPlayer(ctx, input)` だけで検索結果を返せる

### Step 4: Usecase Switch To DB-First

状態: 未着手

目的:

- cache なしで `DB hit -> return / DB miss -> Riot -> DB save -> return` を成立させる

実装要件:

- [backend/internal/usecase/player_search.go](/Users/ellery/dev/koshinankin/feed-gg/backend/internal/usecase/player_search.go) の `Execute` を以下へ変更する
  1. `region` マスタ確認
  2. `repository.FindSavedPlayer`
  3. DB hit なら返す
  4. DB miss なら `riotGateway.SearchPlayerByRiotID`
  5. `repository.SaveFetchedPlayer`
  6. `mapPlayerSearchResult(...)` して返す

完了条件:

- 初回検索は Riot API 経由
- 2回目検索は DB 経由

### Step 5: Cache-Aside

状態: 未着手

目的:

- 検索結果 aggregate をメモリキャッシュに載せる

実装要件:

- `go-cache` を導入する
- `PlayerSearchCache` の本実装を追加する
- `PlayerSearchKey` を使って `cache.Get/Set` する
- まずは positive cache のみ
- `Execute` の流れを以下へ拡張する
  1. `region` マスタ確認
  2. `key := NewPlayerSearchKey(input)`
  3. `cache.Get`
  4. miss の場合だけ `repository.FindSavedPlayer`
  5. DB hit なら `cache.Set`
  6. DB miss なら Riot -> save -> `cache.Set`

完了条件:

- `cache hit -> 即 return`
- `cache miss -> DB or Riot`

### Step 6: Riot API Parallelization

状態: 未着手

目的:

- 検索の応答時間を短縮する

実装要件:

- `account` 取得後に以下を並列化する
  - `summoner`
  - `rank`
  - `match ids`
- `match detail` は `errgroup` 等で制限付き並列にする
- 20件全同時実行ではなく、bounded concurrency にする

完了条件:

- 挙動は変えずに応答時間を改善できる

### Step 7: Normalize And Error Handling

状態: 未着手

目的:

- 入力の一貫性とエラーの扱いを仕上げる

実装要件:

- Go 側で入力正規化する
  - `region` の trim / uppercase
  - `gameName`, `tagLine` の trim
  - 空文字判定
  - rank / participant snapshot の正規化
- Riot 404 / 429 / 5xx の backend エラー方針を決める

完了条件:

- cache key, DB lookup, Riot request の入力が揃う
- frontend が Riot 依存のエラーメッセージに引っ張られない

## Notes For Next Implementation Chat

- `PlayerSearchResult` は当面そのまま使う
  - `cache` と `DB read` の返り値もこれに揃える
- `riot.PlayerProfile` は取得専用の型として扱う
- `update/refresh` 系の経路は今は実装対象外
- `Riot API` を待ちながら transaction を開かない
- `Step 2 -> Step 3 -> Step 4 -> Step 5` の順で進める
