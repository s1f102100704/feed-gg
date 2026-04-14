# TODO

- Normalize Riot API error responses before returning them to the frontend.
  Current implementation forwards upstream error messages almost as-is, which can make the frontend depend on Riot-specific wording and expose implementation details. Map common statuses such as 404/429/5xx to backend-defined error messages.
- Move input normalization from SQL to Go before persisting or querying player/region/rank data.
  `sqlc` queries no longer trim, upper-case, or convert empty strings to NULL. Normalize region names, player names, tag lines, ranks, and optional match participant fields in Go before calling the generated queries.
