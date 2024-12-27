package chessEngine

type Evaluator interface {
	Evaluate(Board) int
}

/*
- Tactical Hueristics

1. 	Mate in 1, 2, 3
2.	Material Difference
3.	Threats

- Positional Heuristics

1.	Knights in center
2.	Bishops on bigger diagonals
3.	Rooks on open files
4. 	KingSafety
*/

type evaluator struct {
}

func (e *evaluator) Evaluate(board *Board, side bool) float32 {
	// return float32(e.MaterialDifference(board))
	return float32(e.MaterialDifference(board)) + 2*e.PieceDevelopment(board, side)
}

func (e *evaluator) MaterialDifference(board *Board) int {
	white := 0
	black := 0
	var i, j int16
	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			piece := board.Get(&Pos{i, j})
			if piece != -1 {
				if piece.Color() {
					black += PIECE_VALUE[piece]
				} else {
					white += PIECE_VALUE[piece]
				}
			}
		}
	}
	return white - black
}

func (e *evaluator) PieceDevelopment(board *Board, side bool) float32 {
	piece_penality := map[Piece]float32{WHITE_ROOK: 0.1, WHITE_BISHOP: 0.4, WHITE_KNIGHT: 0.5, WHITE_QUEEN: 0.1, WHITE_KING: 0.0}
	penality := float32(0)
	if !side {
		if board.Get(&Pos{0, 1}) == WHITE_KNIGHT {
			penality -= piece_penality[WHITE_KNIGHT]
		}
		if board.Get(&Pos{0, 6}) == WHITE_KNIGHT {
			penality -= piece_penality[WHITE_KNIGHT]
		}
		if board.Get(&Pos{0, 2}) == WHITE_BISHOP {
			penality -= piece_penality[WHITE_BISHOP]
		}
		if board.Get(&Pos{0, 5}) == WHITE_BISHOP {
			penality -= piece_penality[WHITE_BISHOP]
		}
		if board.Get(&Pos{0, 0}) == WHITE_ROOK {
			penality -= piece_penality[WHITE_ROOK]
		}
		if board.Get(&Pos{0, 7}) == WHITE_ROOK {
			penality -= piece_penality[WHITE_ROOK]
		}
		if board.Get(&Pos{0, 3}) == WHITE_QUEEN {
			penality -= piece_penality[WHITE_QUEEN]
		}
		if board.Get(&Pos{0, 4}) == WHITE_KING {
			penality -= piece_penality[WHITE_KING]
		}
	} else {
		if board.Get(&Pos{7, 1}) == BLACK_KNIGHT {
			penality += piece_penality[WHITE_KNIGHT]
		}
		if board.Get(&Pos{7, 6}) == BLACK_KNIGHT {
			penality += piece_penality[WHITE_KNIGHT]
		}
		if board.Get(&Pos{7, 2}) == BLACK_BISHOP {
			penality += piece_penality[WHITE_BISHOP]
		}
		if board.Get(&Pos{7, 5}) == BLACK_BISHOP {
			penality += piece_penality[WHITE_BISHOP]
		}
		if board.Get(&Pos{7, 0}) == BLACK_ROOK {
			penality += piece_penality[WHITE_ROOK]
		}
		if board.Get(&Pos{7, 7}) == BLACK_ROOK {
			penality += piece_penality[WHITE_ROOK]
		}
		if board.Get(&Pos{7, 3}) == BLACK_QUEEN {
			penality += piece_penality[WHITE_QUEEN]
		}
		if board.Get(&Pos{7, 4}) == BLACK_KING {
			penality += piece_penality[WHITE_KING]
		}
	}
	return penality
}

var PAWN_TABLE = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	50, 50, 50, 50, 50, 50, 50, 50,
	10, 10, 20, 30, 30, 20, 10, 10,
	5, 5, 10, 25, 25, 10, 5, 5,
	0, 0, 0, 20, 20, 0, 0, 0,
	5, -5, -10, 0, 0, -10, -5, 5,
	5, 10, 10, -20, -20, 10, 10, 5,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var KNIGHTS_TABLE = [64]int{
	-50, -40, -30, -30, -30, -30, -40, -50,
	-40, -20, 0, 0, 0, 0, -20, -40,
	-30, 0, 10, 15, 15, 10, 0, -30,
	-30, 5, 15, 20, 20, 15, 5, -30,
	-30, 0, 15, 20, 20, 15, 0, -30,
	-30, 5, 10, 15, 15, 10, 5, -30,
	-40, -20, 0, 5, 5, 0, -20, -40,
	-50, -40, -30, -30, -30, -30, -40, -50,
}

var BISHOPS_TABLE = [64]int{
	-20, -10, -10, -10, -10, -10, -10, -20,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 5, 10, 10, 5, 0, -10,
	-10, 5, 5, 10, 10, 5, 5, -10,
	-10, 0, 10, 10, 10, 10, 0, -10,
	-10, 10, 10, 10, 10, 10, 10, -10,
	-10, 5, 0, 0, 0, 0, 5, -10,
	-20, -10, -10, -10, -10, -10, -10, -20,
}

var ROOKS_TABLE = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	5, 10, 10, 10, 10, 10, 10, 5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	0, 0, 0, 5, 5, 0, 0, 0,
}
var QUEENS_TABLE = [64]int{
	-20, -10, -10, -5, -5, -10, -10, -20,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 5, 5, 5, 5, 0, -10,
	-5, 0, 5, 5, 5, 5, 0, -5,
	0, 0, 5, 5, 5, 5, 0, -5,
	-10, 5, 5, 5, 5, 5, 0, -10,
	-10, 0, 5, 0, 0, 0, 0, -10,
	-20, -10, -10, -5, -5, -10, -10, -20,
}

var KINGS_TABLE = [64]int{
	-50, -40, -30, -20, -20, -30, -40, -50,
	-30, -20, -10, 0, 0, -10, -20, -30,
	-30, -10, 20, 30, 30, 20, -10, -30,
	-30, -10, 30, 40, 40, 30, -10, -30,
	-30, -10, 30, 40, 40, 30, -10, -30,
	-30, -10, 20, 30, 30, 20, -10, -30,
	-30, -30, 0, 0, 0, 0, -30, -30,
	-50, -30, -30, -30, -30, -30, -30, -50,
}
