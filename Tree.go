package chessEngine

import (
	"math"
	"math/rand/v2"
)

type ChessTree interface {
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
		return n.eval, n.topMoves
	}
}
