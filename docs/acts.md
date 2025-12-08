# Acts

Runs have three acts. Each act is a grid of nodes with edges between rows. Only the current act is visible in the TUI.

Act-specific rules:
- Difficulty filtering: see `docs/constraints.md`.
- Pool sizing: see `docs/constraints.md`.
- Navigation: you commit to a node per row and can only move along reachable edges (no free horizontal moves across the map).
- Goals: each act sets a star target average for its challenges (3★ act 1, 4★ act 2, 5★ act 3).

Advancing between acts resets selection state while keeping the seeded run.
