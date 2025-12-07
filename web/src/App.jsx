import './App.css'
import { generateRun } from './lib/path'
import { useMemo } from 'react'

function App() {
  const { acts, seed } = useMemo(() => generateRun(Date.now()), [])

  return (
    <main className="app">
      <p className="eyebrow">Long Way To The Top</p>
      <h1>Rhythm roguelike — React prototype</h1>
      <p className="lede">
        Hello! This is the starting point for the web client. We&apos;ll add sector maps,
        challenge selection, and rhythm hooks here.
      </p>

      <section className="acts">
        {acts.map((act) => (
          <ActView key={act.index} act={act} />
        ))}
      </section>

      <p className="seed">Seed: {seed}</p>
    </main>
  )
}

function ActView({ act }) {
  // boss at the top: render rows reversed
  const rows = [...act.rows].reverse()
  return (
    <div className="act">
      <h2>Act {act.index}</h2>
      <div className="grid">
        {rows.map((row, rowIdx) => (
          <div className="row" key={`r-${rowIdx}`}>
            {row.map((node) => (
              <div className={`node node-${node.kind}`} key={`c-${rowIdx}-${node.col}`}>
                {node.kind === 'boss' ? 'B' : 'C'}
              </div>
            ))}
          </div>
        ))}
      </div>
    </div>
  )
}

export default App
