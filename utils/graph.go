package utils

import "gonum.org/v1/gonum/graph"

type Node int64

func (n Node) ID() int64 {
	return n.ID()
}

func NewNode(id int) graph.Node {
	return Node(int64(id))
}
