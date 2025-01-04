package chessEngine

import (
	"math"
	"math/rand/v2"
)

type ChessTree interface {
}

var MEMIOZE [10]map[uint64]struct {
	eval  float32
	moves []Move
}

type Node struct {
	evaluator
	board    Board
	children []struct {
		*Node
		Move
	}
	topMoves []Move
	eval     float32
	depth    int
}

func NewNode(board Board, depth int) *Node {
	return &Node{board: board, children: nil, topMoves: nil, depth: depth}
}

func (n *Node) FindChildren(turn bool) {
	var children []struct {
		*Node
		Move
	}
	moves := n.board.GenerateMoves(turn)
	for _, move := range moves {
		chBoard := n.board
		chBoard.Place(move.piece, move.final)
		chBoard.Place(-1, move.init)
		children = append(children, struct {
			*Node
			Move
		}{&Node{board: chBoard, children: nil}, move})
	}
	n.children = children
}

func (n *Node) EvaluateTree(depth int, turn bool) (float32, []Move) {
	value, exists := MEMIOZE[depth][n.board.Hash()]
	if exists {
		return value.eval, value.moves
	}
	if depth == 0 {
		return n.Evaluate(&n.board, turn), []Move{}
	}
	n.FindChildren(turn)
	var childEvals []struct {
		float32
		Move
	}
	for _, child := range n.children {
		eval, _ := child.Node.EvaluateTree(depth-1, !turn)
		childEvals = append(childEvals, struct {
			float32
			Move
		}{eval, child.Move})
	}
	if turn {
		var mn float32 = math.MaxFloat32
		for _, eval := range childEvals {
			if mn > eval.float32 {
				n.topMoves = []Move{eval.Move}
				n.eval = eval.float32
				mn = eval.float32
			} else if mn == eval.float32 && rand.IntN(2) < 1 {
				n.topMoves = []Move{eval.Move}
				n.eval = eval.float32
				mn = eval.float32
			}
		}
		MEMIOZE[depth][n.board.Hash()] = struct {
			eval  float32
			moves []Move
		}{n.eval, n.topMoves}
		return n.eval, n.topMoves
	} else {
		var mx float32 = -math.MaxFloat32
		for _, eval := range childEvals {
			if mx < eval.float32 {
				n.topMoves = []Move{eval.Move}
				n.eval = eval.float32
				mx = eval.float32
			} else if mx == eval.float32 && rand.IntN(2) < 1 {
				n.topMoves = []Move{eval.Move}
				n.eval = eval.float32
				mx = eval.float32
			}
		}
		MEMIOZE[depth][n.board.Hash()] = struct {
			eval  float32
			moves []Move
		}{n.eval, n.topMoves}
		return n.eval, n.topMoves
	}
}

func (n *Node) EvaluateTreeWithPruning(depth int, turn bool, alpha, beta float32) (float32, []Move) {
	if depth == 0 {
		return n.Evaluate(&n.board, turn), []Move{}
	}
	n.FindChildren(turn)
	if turn {
		var mn float32 = math.MaxFloat32
		for _, child := range n.children {
			eval, _ := child.Node.EvaluateTreeWithPruning(depth-1, !turn, alpha, beta)
			if mn > eval {
				n.topMoves = []Move{child.Move}
				n.eval = eval
				mn = eval
			}
			if beta > n.eval {
				beta = n.eval
			}
			if alpha >= beta {
				break
			}
		}
		return n.eval, n.topMoves
	} else {
		var mx float32 = -math.MaxFloat32
		for _, child := range n.children {
			eval, _ := child.Node.EvaluateTreeWithPruning(depth-1, !turn, alpha, beta)
			if mx < eval {
				n.topMoves = []Move{child.Move}
				n.eval = eval
				mx = eval
			}
			if alpha < n.eval {
				alpha = n.eval
			}
			if alpha >= beta {
				break
			}
		}
		return n.eval, n.topMoves
	}
}

func (n *Node) EvaluateTreeConcurrent(depth int, turn bool, response chan struct {
	float32
	Move
}) {
	if depth == 0 {
		response <- struct {
			float32
			Move
		}{n.Evaluate(&n.board, turn), Move{}}
		return
	}
	n.FindChildren(turn)
	// var wg sync.WaitGroup
	var childEvals []chan struct {
		float32
		Move
	}
	var childMoves []Move
	for _, child := range n.children {
		childEvals = append(childEvals, make(chan struct {
			float32
			Move
		}))
		childMoves = append(childMoves, child.Move)
		go child.Node.EvaluateTreeConcurrent(depth-1, !turn, childEvals[len(childEvals)-1])
	}
	if turn {
		var mn float32 = math.MaxFloat32
		for idx, eval := range childEvals {
			ev := <-eval
			move := childMoves[idx]
			if mn > ev.float32 {
				n.topMoves = []Move{move}
				n.eval = ev.float32
				mn = ev.float32
			} else if mn == ev.float32 && rand.IntN(2) < 1 {
				n.topMoves = []Move{move}
				n.eval = ev.float32
				mn = ev.float32
			}
		}
		response <- struct {
			float32
			Move
		}{mn, n.topMoves[0]}
	} else {
		var mx float32 = -math.MaxFloat32
		for idx, eval := range childEvals {
			ev := <-eval
			move := childMoves[idx]
			if mx < ev.float32 {
				n.topMoves = []Move{move}
				n.eval = ev.float32
				mx = ev.float32
			} else if mx == ev.float32 && rand.IntN(2) < 1 {
				n.topMoves = []Move{move}
				n.eval = ev.float32
				mx = ev.float32
			}
		}
		response <- struct {
			float32
			Move
		}{mx, n.topMoves[0]}
	}
}
