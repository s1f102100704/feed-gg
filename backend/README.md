# Backend

## sqlc

`sqlc` generates Go code from SQL files.

- Put schema files in `migrations/`
- Put application queries in `query/`
- Generated code is written to `internal/db/`

Run:

```bash
sqlc generate
```
