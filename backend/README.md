# Backend

## sqlc

`sqlc` generates Go code from SQL files.

- Put schema files in `migrations/`
- Put application queries in `query/`
- Generated code is written to `internal/infrastructure/db/sqlc/`

Run:

```bash
sqlc generate
```

## Lint

Go の lint は `golangci-lint` を Docker 経由で実行する。

```bash
make lint
```

初回実行時は `golangci/golangci-lint` イメージの取得が走る。


/domain
アプリの中心の型と業務ルール
例: Player, Rank, Region, バリデーションやドメイン固有の判定
HTTP/DB/Riot API の都合はなるべく入れない

/usecase
「何をするアプリか」という処理の流れ
例: SearchPlayer, SaveMatchHistory
domain の型を使い、必要な外部依存は interface として定義する
http.Request や SQL の生クエリみたいな外側の詳細は持ち込まない


/adapter
外側との入出力を usecase 用に変換する層
例: HTTP handler、request/response DTO、DB repository のインターフェース定義
GET /api/players/{region}/{gameName}/{tagLine} の path param を usecase input に変換するのはここ


/infrastructure
DB、Riot API、Webフレームワークなど具体技術の実装
例: PostgreSQL repository 実装、Riot API client 実装
/adapter や /usecase で定義した interface を満たす具体実装を置く
