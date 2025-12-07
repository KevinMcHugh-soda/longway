import './App.css'
import { generateRun } from './lib/path'
import { useMemo, useState } from 'react'

function App() {
  const { acts, seed } = useMemo(() => generateRun(Date.now()), [])
  const [currentAct, setCurrentAct] = useState(0)

  return (
    <main className="app">
      <p className="eyebrow">Long Way To The Top</p>
      <h1>Rhythm roguelike — React prototype</h1>
      <p className="lede">
        Hello! This is the starting point for the web client. We&apos;ll add sector maps,
        challenge selection, and rhythm hooks here.
      </p>

      <div className="controls">
        <button onClick={() => setCurrentAct((a) => Math.max(0, a - 1))} disabled={currentAct === 0}>
          Prev Act
        </button>
        <span>
          Act {currentAct + 1} / {acts.length}
        </span>
        <button
          onClick={() => setCurrentAct((a) => Math.min(acts.length - 1, a + 1))}
          disabled={currentAct === acts.length - 1}
        >
          Next Act
        </button>
      </div>

      <section className="acts">
        <ActView act={acts[currentAct]} />
      </section>

      <p className="seed">Seed: {seed}</p>
    </main>
  )
}

function ActView({ act }) {
  // boss at the top: render rows reversed
  const rows = [...act.rows].reverse()
  const height = rows.length * 80

  return (
    <div className="act">
      <h2>Act {act.index}</h2>
      <div className="grid" style={{ height }}>
        {rows.map((row, rowIdx) => (
          <RowView
            key={`r-${rowIdx}`}
            row={row}
            rowIdx={rowIdx}
            totalRows={rows.length}
            nextRow={rowIdx < rows.length - 1 ? rows[rowIdx + 1] : null}
          />
        ))}
      </div>
    </div>
  )
}

function RowView({ row, rowIdx, totalRows, nextRow }) {
  return (
    <div className="row" style={{ top: rowIdx * 80 }}>
      {row.map((node) => (
        <NodeView key={`n-${rowIdx}-${node.col}`} node={node} rowIdx={rowIdx} nextRow={nextRow} />
      ))}
    </div>
  )
}

function NodeView({ node, rowIdx, nextRow }) {
  const x = node.col * 80

  const edges = (nextRow || []).map((_, idx) => idx)

  return (
    <div className="node-wrapper" style={{ left: x }}>
      {nextRow &&
        (node.edges || []).map((target) => {
          const targetCol = edges[target] ?? target
          const dx = (targetCol - node.col) * 80
          const dy = 80
          const angle = Math.atan2(dy, dx) * (180 / Math.PI)
          const length = Math.sqrt(dx * dx + dy * dy)
          return (
            <div
              key={`edge-${target}`}
              className="edge"
              style={{
                width: `${length}px`,
                transform: `translate(20px, 20px) rotate(${angle}deg)`,
              }}
            />
          )
        })}
      <div className={`node node-${node.kind}`}>{node.kind === 'boss' ? 'B' : 'C'}</div>
    </div>
  )
}

export default App
