package main

type node struct {
	col       int
	edges     []int // indices into the next row
	kind      nodeKind
	challenge *challenge
}

type act struct {
	index int
	rows  [][]node
}

type nodeKind int

const (
	nodeUnknown nodeKind = iota
	nodeChallenge
	nodeShop
	nodeBoss
)

const (
	totalActs             = 3
	rowsPerAct            = 7
	minNodesPerRow        = 2
	maxNodesPerRow        = 3
	colSpacing            = 4
	challengeSongListSize = 12
)

type nodeRun struct {
	col   int
	songs []song
	stars []int
}
