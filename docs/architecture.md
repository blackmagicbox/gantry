## Goal
A block-stacking game where you are a Gantry Crane operator that needs to load cargo ships with puzzeling container shapes, fully playable offline as a single-player game, with online multiplayer limited to duels and tournaments. Client: Kotlin Multiplatform + Compose Multiplatform, Android-first, sharing all game logic across platforms — the only per-platform variation is login (Google/Apple/email) and, later, optional Game Center/Play Games integration. Backend: Go microservices on GKE Autopilot, handling matchmaking, duel/tournament session sync, and leaderboards.

## Services
### `matchmaking-service`
**Owns**: pairing players (queue-based matches and direct rematches) and triggering duel-service to start a match.

**Does not own**: live match state, score tracking, disconnect/quit handling, or outcome recording during a duel (all duel-service) — matchmaking's job ends the moment a duel is triggered, whether the pairing came from the queue or a rematch request.

### `duel-service`
**Owns**: live match state, score tracking, detecting game-end conditions (topout, quit, timeout), declaring the winner, notifying players, and calling leaderboard-service to record the result.

**Does not own**: pairing players, rematch requests (matchmaking-service), persisting historical scores, computing rankings, or displaying player standings (leaderboard-service).

### `leaderboard-service`
**Owns**: computing and storing ELO/rating changes from reported duel results, tracking scores within a ladder/season window, ranking by rating + windowed points, and displaying standings (including tournament framing as a season window).

**Does not own**: deciding who won a duel or when it ends (duel-service) leaderboard-service only receives the final result and score, it doesn't observe or judge the match itself.

### `auth-service`
**Owns**: account existence and credentials (email/provider ID), issuing and validating access tokens that other services use to identify which player is making a request.

**Does not own**: player profile customization (display name, avatar, etc.) — none exists yet; that's a future user-service if/when it's needed. Also does not own pairing, matches, or any game state, and does not handle service-to-service trust between backend services (that's an infrastructure/network-policy concern, not application logic).
