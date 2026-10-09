# Failure modes

One line per call/mechanism: what breaks, and what happens

## `matchmaking-service` -> `duel-service`

TriggerDuel call (`matchmaking-service` -> `duel-service`): if it fails or times out, retry up to 5 times with increasing backoff; after the fifth failure, surface "unable to start duel, try again later" to the client

TriggerDuel retries: idempotency_key tracks a status (in-progress/success/failed) per attempt. A retry with the same key returns the same result if succeed instead of creating a new duel.

## `matchmaking-service` - Known Limitations

Client authentication: the authentication interceptor requires an `x-player-id`
gRPC metadata header, rejects missing or empty values with `Unauthenticated`,
and exposes the value to handlers through `PlayerIDFromContext`. The header
is client-supplied and not verified, so any client can impersonate any
player. This is acceptable for local development only; the service must not
be exposed beyond it. Once auth-service can issue tokens, the interceptor
will validate a token and derive the player ID from it. Handlers will not
change. A ticket must then be bound to the player who created it, so that
GetQueueStatus rejects a ticket_id that belongs to someone else.
JoinQueueRequest is empty for that reason: identity never comes from the
request body.

Queue persistence: tickets live in an in-memory map only. If
matchmaking-service restarts, all waiting tickets are lost and players must
rejoin. Tickets are also never removed, so the map grows without bound;
expiry (STATUS_EXPIRED) is planned but not implemented. Acceptable at
current scale (single instance, pre-launch); revisit with expiry and
persistent storage (Redis/DB) before running multiple replicas or handling
real traffic.

## `duel-service`

Duel timeout checkpoint: if a client doesn't respond within the grace
window, treat that player as having quit (forfeit, no score comparison) —
same rule as an explicit quit. (Reconnection/re-matchmaking after a
dropped client is out of scope for now.)

## `duel-service` → `leaderboard-service`

UpdateLeaderboard retries: same idempotency_key pattern as TriggerDuel — a
retry with the same match_id/key returns the already-recorded rating
change instead of double-counting.

## `duel-service` - Known Limitations

Idempotency key persistence: idempotency_key state (TriggerDuel) lives in
an in-memory map only. If duel-service restarts between creating a match
and returning the response, a subsequent retry with the same key is
treated as new and creates a duplicate match. Acceptable at current scale
(single instance, pre-launch); revisit with persistent storage (Redis/DB)
before running multiple replicas or handling real traffic.
