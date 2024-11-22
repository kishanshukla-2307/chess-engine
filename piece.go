package chessEngine

type Piece int16

func (p Piece) Color() bool {
	return p > 5
}

func (p Piece) IsKing() bool {
	return p == 0 || p == 6
}

func (p Piece) IsPawn() bool {
	return p == 5 || p == 11
}

// func (p Piece) CanReach(crr Position, dest Position) bool {
// 	switch {
// 	case p == WHITE_KING || p == BLACK_KING:
// 		return dest.Sub(crr).GetRank()|dest.Sub(crr).GetRank() == 1
// 	case p == WHITE_KING || p == BLACK_KING:
// 		return dest.Sub(crr).GetRank()|dest.Sub(crr).GetRank() == 1
// 	case p == WHITE_KING || p == BLACK_KING:
// 		return dest.Sub(crr).GetRank()|dest.Sub(crr).GetRank() == 1
// 	case p == WHITE_KING || p == BLACK_KING:
// 		return dest.Sub(crr).GetRank()|dest.Sub(crr).GetRank() == 1
// 	case p == WHITE_KING || p == BLACK_KING:
// 		return dest.Sub(crr).GetRank()|dest.Sub(crr).GetRank() == 1
// 	}
// }

const (
	WHITE_KING   Piece = 0
	WHITE_QUEEN  Piece = 1
	WHITE_ROOK   Piece = 2
	WHITE_BISHOP Piece = 3
	WHITE_KNIGHT Piece = 4
	WHITE_PAWN   Piece = 5
	BLACK_KING   Piece = 6
	BLACK_QUEEN  Piece = 7
	BLACK_ROOK   Piece = 8
	BLACK_BISHOP Piece = 9
	BLACK_KNIGHT Piece = 10
	BLACK_PAWN   Piece = 11
)

var (
	pieceUnicodes = []string{"♔", "♕", "♖", "♗", "♘", "♙", "♚", "♛", "♜", "♝", "♞", "♟"}
)
