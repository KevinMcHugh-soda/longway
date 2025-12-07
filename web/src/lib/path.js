import { songs as catalog } from '../data/songs'

const totalActs = 3
const rowsPerAct = 8
const minNodesPerRow = 2
const maxNodesPerRow = 5

const poolBounds = {
  1: { min: 9, max: 12 },
  2: { min: 6, max: 9 },
  3: { min: 3, max: 5 },
}

export const nodeKinds = {
  challenge: 'challenge',
  boss: 'boss',
}

export function generateRun(seed = Date.now()) {
  const rng = mulberry32(seed)
  const acts = []
  for (let i = 0; i < totalActs; i++) {
    acts.push(generateAct(i + 1, rng))
  }
  return { acts, seed }
}

function generateAct(index, rng) {
  const filteredSongs = applyActDifficultyConstraints(index, catalog)
  const poolSize = pickPoolSize(index, filteredSongs.length, rng)
  const rows = []

  for (let row = 0; row < rowsPerAct; row++) {
    let count = minNodesPerRow + rngInt(rng, maxNodesPerRow - minNodesPerRow + 1)
    if (row === rowsPerAct - 1) {
      count = 1 // boss
    }

    const nodes = []
    for (let col = 0; col < count; col++) {
      const isBoss = row === rowsPerAct - 1
      nodes.push({
        col,
        kind: isBoss ? nodeKinds.boss : nodeKinds.challenge,
        challenge: isBoss ? bossChallenge() : challenge(filteredSongs, poolSize, rng),
        edges: [],
      })
    }

    if (row > 0) {
      connectRows(rows[row - 1], nodes, rng)
    }

    rows.push(nodes)
  }

  return { index, rows }
}

function connectRows(prev, next, rng) {
  if (!prev.length || !next.length) return

  // ensure every next node has an inbound edge
  next.forEach((_, nextIdx) => {
    const src = rngInt(rng, prev.length)
    prev[src].edges.push(nextIdx)
  })

  // add extra edges
  prev.forEach((node) => {
    const targetCount = Math.min(1 + rngInt(rng, 2), next.length)
    const picks = pickDistinct(next.length, targetCount, rng)
    node.edges.push(...picks)
  })
}

function challenge(pool, poolSize, rng) {
  const songs = sample(pool, poolSize, rng)
  return {
    name: 'Challenge',
    summary: `Pick any 3 of these ${songs.length} tracks.`,
    songs,
  }
}

function bossChallenge() {
  const boss = catalog.find((s) => s.title === 'Bohemian Rhapsody') ?? catalog[0]
  return {
    name: 'Boss',
    summary: 'Final showdown: Bohemian Rhapsody.',
    songs: [boss],
  }
}

function applyActDifficultyConstraints(actIndex, songs) {
  return songs.filter((s) => {
    if (actIndex === 1) return s.difficulty <= 3
    if (actIndex === 2) return s.difficulty <= 5
    return s.difficulty >= 3
  })
}

function pickPoolSize(actIndex, available, rng) {
  const bounds = poolBounds[actIndex] ?? poolBounds[3]
  const min = Math.min(bounds.min, available)
  const max = Math.min(bounds.max, available)
  if (min >= max) return min
  return min + rngInt(rng, max - min + 1)
}

function sample(pool, count, rng) {
  if (pool.length <= count) return [...pool]
  const indices = pickDistinct(pool.length, count, rng)
  return indices.map((i) => pool[i])
}

function pickDistinct(size, count, rng) {
  const seen = new Set()
  const picks = []
  while (picks.length < count) {
    const v = rngInt(rng, size)
    if (seen.has(v)) continue
    seen.add(v)
    picks.push(v)
  }
  return picks
}

function rngInt(rng, maxExclusive) {
  return Math.floor(rng() * maxExclusive)
}

// deterministic PRNG
function mulberry32(seed) {
  let t = seed + 0x6d2b79f5
  return function () {
    t = Math.imul(t ^ (t >>> 15), t | 1)
    t ^= t + Math.imul(t ^ (t >>> 7), t | 61)
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296
  }
}
