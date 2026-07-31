## Goal
A from-scratch Tetris rebuild, fully playable offline as a single-player game, with online multiplayer limited to duels and tournaments. Client: Kotlin Multiplatform + Compose Multiplatform, Android-first, sharing all game logic across platforms — the only per-platform variation is login (Google/Apple/email) and, later, optional Game Center/Play Games integration. Backend: Go microservices on GKE Autopilot, handling matchmaking, duel/tournament session sync, and leaderboards.

## Services
