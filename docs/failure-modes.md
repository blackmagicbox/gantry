# Failure modes

One line per call/mechanism: what breaks, and what happens

## `matchmaking-service` -> `duel-service`

TriggerDuel call (`matchmaking-service` -> `duel-service`): if it fails or times out, retry up to 5 times with increasing backoff; after the fifth failure, surface "unable to start duel, try again later" to the client

TriggerDuel retries: idempotency_key tracks a status (in-progress/success/failed) per attempt. A retry with the same key returns the same result if succeed instead of creating a new duel.

## `duel-service`

Duel timeout checkpoint: if a client doesn't respond within the grace
window, treat that player as having quit (forfeit, no score comparison) —
same rule as an explicit quit. (Reconnection/re-matchmaking after a
dropped client is out of scope for now.)

## `duel-service` → `leaderboard-service`

UpdateLeaderboard retries: same idempotency_key pattern as TriggerDuel — a
retry with the same match_id/key returns the already-recorded rating
change instead of double-counting.
