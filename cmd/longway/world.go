package main

import "math/rand"

func generateRun(seed int64, songs []song) []act {
	rng := rand.New(rand.NewSource(seed))
	acts := make([]act, totalActs)
	for i := 0; i < totalActs; i++ {
		acts[i] = generateAct(i+1, rng, songs)
	}
	return acts
}

func generateAct(index int, rng *rand.Rand, songs []song) act {
	actSongs := applyActDifficultyConstraints(index, songs)
	poolSize := pickPoolSize(index, len(actSongs), rng)
	shopRows := pickShopRows(rng)
	rows := make([][]node, rowsPerAct)
	for row := 0; row < rowsPerAct; row++ {
		count := minNodesPerRow + rng.Intn(maxNodesPerRow-minNodesPerRow+1)
		if row == rowsPerAct-1 || shopRows[row] {
			count = 1 // boss or shop
		}
		nodes := make([]node, count)
		for i := range nodes {
			kind := nodeChallenge
			if row == rowsPerAct-1 {
				kind = nodeBoss
			} else if shopRows[row] {
				kind = nodeShop
			}
			nodes[i] = node{
				col:  i,
				kind: kind,
			}
			if kind == nodeChallenge {
				nodes[i].challenge = newChallenge(actSongs, rng, poolSize)
			} else if kind == nodeBoss {
				nodes[i].challenge = bossChallenge(songs)
			}
		}
		if row > 0 {
			connectRows(rows[row-1], nodes, rng)
		}
		rows[row] = nodes
	}

	return act{
		index: index,
		rows:  rows,
	}
}

func pickShopRows(rng *rand.Rand) map[int]bool {
	shopCount := 2
	candidates := []int{}
	for r := 1; r < rowsPerAct-1; r++ { // avoid first and last rows
		candidates = append(candidates, r)
	}
	selected := map[int]bool{}
	attempts := 0
	for len(selected) < shopCount && attempts < 100 {
		r := candidates[rng.Intn(len(candidates))]
		// no back-to-back
		if selected[r-1] || selected[r+1] {
			attempts++
			continue
		}
		selected[r] = true
	}
	// ensure we have the desired count by spacing if needed
	for len(selected) < shopCount && len(candidates) > 0 {
		for _, r := range candidates {
			if selected[r] || selected[r-1] || selected[r+1] {
				continue
			}
			selected[r] = true
			if len(selected) >= shopCount {
				break
			}
		}
	}
	return selected
}

func connectRows(prev []node, next []node, rng *rand.Rand) {
	incoming := make([]int, len(next))

	// ensure every next node has an inbound edge by building a spanning pass first
	for j := range next {
		src := rng.Intn(len(prev))
		prev[src].edges = append(prev[src].edges, j)
		incoming[j]++
	}

	// add extra edges for branching
	for i := range prev {
		targets := pickTargets(len(next), rng)
		for _, t := range targets {
			prev[i].edges = append(prev[i].edges, t)
			incoming[t]++
		}
	}
}

func pickTargets(nextCount int, rng *rand.Rand) []int {
	if nextCount <= 0 {
		return nil
	}
	maxTargets := 1
	if nextCount > 1 {
		maxTargets = 2
	}
	targetCount := 1 + rng.Intn(maxTargets) // up to maxTargets unique targets
	targets := make([]int, 0, targetCount)
	seen := make(map[int]struct{})
	for len(targets) < targetCount {
		t := rng.Intn(nextCount)
		if _, ok := seen[t]; ok {
			continue
		}
		seen[t] = struct{}{}
		targets = append(targets, t)
	}
	return targets
}

func applyActDifficultyConstraints(actIndex int, songs []song) []song {
	if len(songs) == 0 {
		return songs
	}

	filtered := make([]song, 0, len(songs))
	for _, s := range songs {
		d := clampDifficulty(s.difficulty)

		switch actIndex {
		case 1:
			if d <= 3 {
				filtered = append(filtered, s)
			}
		case 2:
			if d <= 5 {
				filtered = append(filtered, s)
			}
		default:
			if d >= 3 {
				filtered = append(filtered, s)
			}
		}
	}

	if len(filtered) == 0 {
		return songs
	}
	return filtered
}

func pickPoolSize(actIndex int, available int, rng *rand.Rand) int {
	minSize, maxSize := poolBoundsForAct(actIndex)
	if available < minSize {
		return available
	}
	if available < maxSize {
		maxSize = available
	}
	if minSize == maxSize {
		return minSize
	}
	return minSize + rng.Intn(maxSize-minSize+1)
}

func poolBoundsForAct(actIndex int) (int, int) {
	switch actIndex {
	case 1:
		return 9, 12
	case 2:
		return 6, 9
	default:
		return 3, 5
	}
}
