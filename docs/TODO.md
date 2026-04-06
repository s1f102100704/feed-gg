# TODO

- Normalize Riot API error responses before returning them to the frontend.
  Current implementation forwards upstream error messages almost as-is, which can make the frontend depend on Riot-specific wording and expose implementation details. Map common statuses such as 404/429/5xx to backend-defined error messages.
