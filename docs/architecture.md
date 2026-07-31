## Goal
A from-scratch block-stacking game re-imagined, fully playable offline as a single-player game, with online multiplayer limited to duels and tournaments. Client: Kotlin Multiplatform + Compose Multiplatform, Android-first, sharing all game logic across platforms — the only per-platform variation is login (Google/Apple/email) and, later, optional Game Center/Play Games integration. Backend: Go microservices on GKE Autopilot, handling matchmaking, duel/tournament session sync, and leaderboards.

## Services
### matchmaking-service
Service to select players and create a match
### duel-service
A service to pair players and handle the duel
### leaderboard-service
A service to store the scores and rank the players, that can also be used to track scores in a tournament
### auth-service
A service to authenticate players.
