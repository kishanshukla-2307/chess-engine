package chessEngine

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"

	"github.com/kishanshukla-2307/chess-engine/utils"
)

type Engine interface {
	IsLegal(Piece, Position, Position) (bool, error)
	MakeMove(Piece, Position, Position) error
}

// var _ Engine = (*NoobEngine)(nil)

type NoobEngine struct {
	board     Board
	chess960  bool
	enPassant bool
}

func NewNoobEngine(chess960 bool) (*NoobEngine, error) {
	var board Board
	err := board.InitializeBoard(chess960)
	if err != nil {
		return nil, err
	}
	return &NoobEngine{board: board,
		chess960: chess960}, nil
}

func (ne *NoobEngine) Run() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		if !ne.board.turn {
			tree := NewNode(ne.board, 4)
			eval, moves := tree.EvaluateTree(4, false)
			fmt.Println(eval)
			move := moves[0]
			err := ne.board.MakeMove(move.piece, move.init, move.final)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("Enter Black's Move: ")
			moveStr, _ := reader.ReadString('\n')
			move := strings.Split(moveStr, " ")
			piece, err := ne.PieceFromNotation(move[0])
			if err != nil {
				return err
			}
			init, err := ne.PositionFromNotation(move[1])
			if err != nil {
				return err
			}
			final, err := ne.PositionFromNotation(move[2])
			if err != nil {
				return err
			}
			err = ne.board.MakeMove(piece, init, final)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func (ne *NoobEngine) PieceFromNotation(n string) (Piece, error) {
	if ne.board.turn {
		switch n {
		case "K":
			return BLACK_KING, nil
		case "Q":
			return BLACK_QUEEN, nil
		case "R":
			return BLACK_ROOK, nil
		case "B":
			return BLACK_BISHOP, nil
		case "N":
			return BLACK_KNIGHT, nil
		case "P":
			return BLACK_PAWN, nil
		}
	} else {
		switch n {
		case "K":
			return WHITE_KING, nil
		case "Q":
			return WHITE_QUEEN, nil
		case "R":
			return WHITE_ROOK, nil
		case "B":
			return WHITE_BISHOP, nil
		case "N":
			return WHITE_KNIGHT, nil
		case "P":
			return WHITE_PAWN, nil
		}
	}
	return -1, errors.New("Unsupported notation")
}

func (ne *NoobEngine) PositionFromNotation(pos string) (Position, error) {
	file := (int16)(pos[0] - 'a')
	rank, err := strconv.Atoi(string(pos[1]))
	if err != nil {
		return &Pos{}, nil
	}
	rank--
	var res Position = &Pos{(int16)(rank), file}
	return res, nil
}

// func (ne *NoobEngine) IsLegal(p Piece, init Position, final Position) (bool, error) {
// 	if init.Equal(final) {
// 		return false, errors.New("init and final are same")
// 	}

// 	// check if final pos have same color piece
// 	if ne.board.Get(final) != -1 {
// 		if p.Color() == ne.board.Get(final).Color() {
// 			return false, errors.New("attacking same color")
// 		}
// 	}

// 	if p.IsKing() {
// 		if !KingAttack(&ne.board, p, init, final) {
// 			return false, errors.New("King cant move like that")
// 		}
// 		if ne.IsAttackedBySide(final, !p.Color()) {
// 			return false, errors.New("King can't move into check")
// 		}
// 	} else if p.IsPawn() {
// 		if !PawnMove(&ne.board, p, init, final) {
// 			return false, errors.New("Pawn can't move there")
// 		}
// 	} else {
// 		if !ne.IsAttackedByPieceWithPos(p)(&ne.board, p, init, final) {
// 			return false, errors.New("piece can't go there")
// 		}
// 	}

// 	if ne.inCheck {

// 		tarPiece := ne.board.Get(final)
// 		ne.board.MakeMove(p, init, final)
// 		var kingPos Position
// 		if p.Color() {
// 			kingPos = ne.board.GetPiecePositions(BLACK_KING)[0]
// 		} else {
// 			kingPos = ne.board.GetPiecePositions(WHITE_KING)[0]
// 		}
// 		if ne.IsAttackedBySide(kingPos, !p.Color()) {
// 			return false, errors.New("in check!")
// 		}
// 		ne.board.Place(tarPiece, final)
// 		ne.board.Place(p, init)
// 	}
// 	return true, nil
// }

// // checks if pos is attacked by side
// func (ne *NoobEngine) IsAttackedBySide(pos Position, side bool) bool {
// 	startingPiece := WHITE_KING
// 	if side {
// 		startingPiece = BLACK_KING
// 	}
// 	for piece := startingPiece; piece < startingPiece+6; piece += 1 {
// 		if ne.IsAttackedByPiece(piece, pos) {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (ne *NoobEngine) IsAttackedByPiece(piece Piece, tar Position) bool {
// 	var i, j int16
// 	for i = 0; i < 8; i++ {
// 		for j = 0; j < 8; j++ {
// 			if ne.board[i][j] == piece {
// 				if i == tar.GetRank() && j == tar.GetFile() {
// 					continue
// 				}
// 				var pos Position = &Pos{i, j}
// 				if ne.IsAttackedByPieceWithPos(piece)(&ne.board, piece, pos, tar) {
// 					return true
// 				}
// 			}
// 		}
// 	}
// 	return false
// }

// func (ne *NoobEngine) IsAttackedByPieceWithPos(piece Piece) func(*Board, Piece, Position, Position) bool {
// 	switch piece {
// 	case WHITE_KING | BLACK_KING:
// 		return KingAttack
// 	case WHITE_QUEEN | BLACK_QUEEN:
// 		return QueenAttack
// 	case WHITE_ROOK | BLACK_ROOK:
// 		return RookAttack
// 	case WHITE_BISHOP | BLACK_BISHOP:
// 		return BishopAttack
// 	case WHITE_KNIGHT | BLACK_KNIGHT:
// 		return KnightAttack
// 	case WHITE_PAWN | BLACK_PAWN:
// 		return PawnAttack
// 	}
// 	return nil
// }

func (ne *NoobEngine) MakeMoveByNotation(move string) error {
	return nil
}

// func (ne *NoobEngine) MakeMove(p Piece, init Position, final Position) error {
// 	if legal, err := ne.board.IsLegal(p, init, final); !legal {
// 		return errors.New("illegal move: " + err.Error())
// 	}
// 	ne.board.MakeMove(p, init, final)
// 	ne.turn = !ne.turn
// 	oppositionKing := WHITE_KING
// 	if !p.Color() {
// 		oppositionKing = BLACK_KING
// 	}
// 	pos := ne.board.GetPiecePositions(oppositionKing)[0]
// 	if ne.board.IsAttackedBySide(pos, p.Color()) {
// 		ne.inCheck = true
// 	} else {
// 		ne.inCheck = false
// 	}
// 	ne.PrintBoard()
// 	return nil
// }

// func (ne *NoobEngine) PrintBoard() {
// 	fmt.Printf("xx")
// 	for i := 7; i >= 0; i-- {
// 		fmt.Printf("xxx")
// 	}
// 	fmt.Println()
// 	for i := 7; i >= 0; i-- {
// 		fmt.Printf("X")
// 		for j := 0; j < 8; j++ {
// 			if ne.board.squares[i][j] == -1 {
// 				fmt.Printf(" - ")
// 			} else {
// 				fmt.Printf(" " + pieceUnicodes[ne.board.squares[i][j]] + " ")
// 			}
// 		}
// 		fmt.Printf("X")
// 		fmt.Printf("\n")
// 	}
// 	fmt.Printf("xx")
// 	for i := 7; i >= 0; i-- {
// 		fmt.Printf("xxx")
// 	}
// 	fmt.Println()
// 	if ne.inCheck {
// 		if ne.turn {
// 			fmt.Printf("Black ")
// 		} else {
// 			fmt.Printf("White ")
// 		}
// 		fmt.Printf("in check!")
// 	}
// 	fmt.Println()
// }

// func (ne *NoobEngine) GenerateMoves(side bool) []Move {
// 	var moves []Move
// 	if side {
// 		for _, piece := range BLACK_PIECES {
// 			positions := ne.board.GetPiecePositions(piece)
// 			for _, pos := range positions {
// 				moves = append(moves, ne.board.GenerateMovesForPiece(piece)(&ne.board, pos)...)
// 			}
// 		}
// 	} else {
// 		for _, piece := range WHITE_PIECES {
// 			positions := ne.board.GetPiecePositions(piece)
// 			for _, pos := range positions {
// 				moves = append(moves, ne.board.GenerateMovesForPiece(piece)(&ne.board, pos)...)
// 			}
// 		}
// 	}
// 	var legalMoves []Move
// 	for _, move := range moves {
// 		if legal, err := ne.board.IsLegal(move.piece, move.init, move.final); (err == nil) && legal {
// 			legalMoves = append(legalMoves, move)
// 		}
// 	}
// 	return legalMoves
// }

func (ne *NoobEngine) RandomMove(side bool) Move {
	moves := ne.board.GenerateMoves(side)
	fmt.Println(len(moves))
	idx := rand.IntN(len(moves))
	return moves[idx]
}

func (ne *NoobEngine) BestMove(side bool) Move {
	return Move{}
}

func KingAttack(board *Board, piece Piece, orig Position, tar Position) bool {
	diff := tar.Sub(orig)
	return utils.Abs(diff.GetRank())|utils.Abs(diff.GetFile()) == 1
}

func QueenAttack(board *Board, piece Piece, orig Position, tar Position) bool {
	return RookAttack(board, piece, orig, tar) || BishopAttack(board, piece, orig, tar)
}

func RookAttack(board *Board, piece Piece, orig Position, tar Position) bool {
	if (tar.GetRank() == orig.GetRank()) || (tar.GetFile() == orig.GetFile()) {
		if tar.GetRank() == orig.GetRank() {
			start := min(tar.GetFile(), orig.GetFile()) + 1
			for ; start < max(tar.GetFile(), orig.GetFile()); start++ {
				if board.Get(&Pos{orig.GetRank(), start}) != -1 {
					return false
				}
			}
		} else {
			start := min(tar.GetRank(), orig.GetRank()) + 1
			for ; start < max(tar.GetRank(), orig.GetRank()); start++ {
				if board.Get(&Pos{start, orig.GetFile()}) != -1 {
					return false
				}
			}
		}
		return true
	}
	return false
}

func BishopAttack(board *Board, piece Piece, orig Position, tar Position) bool {
	diff := tar.Sub(orig)
	if utils.Abs(diff.GetRank()) == utils.Abs(diff.GetFile()) {
		var rank_ch int16 = -1
		if diff.GetRank() > 0 {
			rank_ch = 1
		}
		var file_ch int16 = -1
		if diff.GetFile() > 0 {
			file_ch = 1
		}

		rank := orig.GetRank() + rank_ch
		file := orig.GetFile() + file_ch

		for rank != tar.GetRank() && file != tar.GetFile() {
			if board.squares[rank][file] != -1 {
				return false
			}
			rank += rank_ch
			file += file_ch
		}
		return true
	}
	return false
}

func KnightAttack(board *Board, piece Piece, orig Position, tar Position) bool {
	diff := tar.Sub(orig)
	return (utils.Abs(diff.GetRank()) == 2 && utils.Abs(diff.GetFile()) == 1) || (utils.Abs(diff.GetRank()) == 1 && utils.Abs(diff.GetFile()) == 2)
}

func PawnAttack(board *Board, piece Piece, orig Position, tar Position) bool {
	diff := tar.Sub(orig)
	if !piece.Color() {
		if (diff.GetRank() == 1 && utils.Abs(diff.GetFile()) == 1) && !board.IsEmpty(tar) && board.Get(tar).Color() != piece.Color() {
			return true
		}
	} else {
		if (diff.GetRank() == -1 && utils.Abs(diff.GetFile()) == 1) && !board.IsEmpty(tar) && board.Get(tar).Color() != piece.Color() {
			return true
		}
	}
	return false
}

func PawnMove(board *Board, piece Piece, orig Position, tar Position) bool {
	diff := tar.Sub(orig)
	if !piece.Color() {
		if (diff.GetRank() == 1 && utils.Abs(diff.GetFile()) == 1) && !board.IsEmpty(tar) && board.Get(tar).Color() != piece.Color() {
			return true
		}
		if (diff.GetRank() == 1 && diff.GetFile() == 0) && board.IsEmpty(tar) {
			return true
		}
		if (diff.GetRank() == 2 && diff.GetFile() == 0) && board.IsEmpty(tar) && board.IsEmpty(tar.Sub(&Pos{1, 0})) {
			return true
		}
	} else {
		if (diff.GetRank() == -1 && utils.Abs(diff.GetFile()) == 1) && !board.IsEmpty(tar) && board.Get(tar).Color() != piece.Color() {
			return true
		}
		if (diff.GetRank() == -1 && diff.GetFile() == 0) && board.IsEmpty(tar) {
			return true
		}
		if (diff.GetRank() == -2 && diff.GetFile() == 0) && board.IsEmpty(tar) && board.IsEmpty(tar.Sub(&Pos{-1, 0})) {
			return true
		}
	}
	return false
}
