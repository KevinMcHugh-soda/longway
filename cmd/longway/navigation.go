package main

func initAllowed(a act) []int {
	cols := make([]int, len(a.rows[0]))
	for i := range cols {
		cols[i] = i
	}
	return cols
}

func (m *model) setAllowedForRow(row int) {
	if row == 0 {
		m.allowed = initAllowed(m.acts[m.currentAct])
		m.allowedIdx = 0
		m.cursorCol = m.allowed[0]
		return
	}

	prevCol, ok := m.committed[row-1]
	if !ok {
		m.allowed = initAllowed(m.acts[m.currentAct])
		m.allowedIdx = 0
		m.cursorCol = m.allowed[0]
		return
	}

	prevRow := m.acts[m.currentAct].rows[row-1]
	if prevCol >= len(prevRow) {
		prevCol = len(prevRow) - 1
	}
	edges := prevRow[prevCol].edges
	if len(edges) == 0 {
		m.allowed = initAllowed(m.acts[m.currentAct])
		m.allowedIdx = 0
		m.cursorCol = m.allowed[0]
		return
	}

	m.allowed = make([]int, len(edges))
	copy(m.allowed, edges)
	m.allowedIdx = 0
	m.cursorCol = m.allowed[0]
}

func (m *model) moveHorizontal(delta int) {
	if m.awaitStars {
		return
	}
	if len(m.allowed) == 0 {
		return
	}
	m.allowedIdx += delta
	if m.allowedIdx < 0 {
		m.allowedIdx = 0
	} else if m.allowedIdx >= len(m.allowed) {
		m.allowedIdx = len(m.allowed) - 1
	}
	m.cursorCol = m.allowed[m.allowedIdx]
}

func (m *model) commitSelection() {
	if m.awaitStars {
		return
	}
	m.committed[m.cursorRow] = m.cursorCol
	m.awaitStars = true
	m.starInput = ""
}

func (m *model) submitStars() {
	if !m.awaitStars {
		return
	}
	val := clampDifficulty(parseDifficulty(m.starInput))
	m.stars[m.cursorRow] = val
	m.awaitStars = false

	if m.cursorRow < len(m.acts[m.currentAct].rows)-1 {
		m.cursorRow++
		m.setAllowedForRow(m.cursorRow)
		m.cursorCol = m.allowed[m.allowedIdx]
	}
}
