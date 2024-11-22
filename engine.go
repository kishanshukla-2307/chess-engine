package chessEngine

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kishanshukla-2307/chess-engine/utils"
)

type Engine interface {
	IsLegal(Piece, Position, Position) (bool, error)
	MakeMove(Piece, Position, Position) error
}

var _ Engine = (*NoobEngine)(nil)

type NoobEngine struct {
	board             Board
	turn              bool
	chess960          bool
	inCheck           bool
	enPassant         bool
	isAttackedByPiece map[Piece]func(*Board, Piece, Position, Position) bool
}

func NewNoobEngine(chess960 bool) (*NoobEngine, error) {
	var board Board
	err := InitializeBoard(chess960, &board)
	if err != nil {
		return nil, err
	}
	turn := false
	isAttackedByPiece := make(map[Piece]func(*Board, Piece, Position, Position) bool)
	isAttackedByPiece[0] = KingAttack
	isAttackedByPiece[6] = KingAttack
	isAttackedByPiece[1] = QueenAttack
	isAttackedByPiece[7] = QueenAttack
	isAttackedByPiece[2] = RookAttack
	isAttackedByPiece[8] = RookAttack
	isAttackedByPiece[3] = BishopAttack
	isAttackedByPiece[9] = BishopAttack
	isAttackedByPiece[4] = KnightAttack
	isAttackedByPiece[10] = KnightAttack
	isAttackedByPiece[5] = PawnAttack
	isAttackedByPiece[11] = PawnAttack
	return &NoobEngine{board: board,
		turn:              turn,
		isAttackedByPiece: isAttackedByPiece,
		chess960:          chess960}, nil
}

func InitializeBoard(chess960 bool, board *Board) error {
	if chess960 {
		return errors.New("chess 960 not supported yet")
	}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			board[i][j] = -1
		}
	}
	board[0] = [8]Piece{WHITE_ROOK, WHITE_KNIGHT, WHITE_BISHOP, WHITE_QUEEN, WHITE_KING, WHITE_BISHOP, WHITE_KNIGHT, WHITE_ROOK}
	board[7] = [8]Piece{BLACK_ROOK, BLACK_KNIGHT, BLACK_BISHOP, BLACK_QUEEN, BLACK_KING, BLACK_BISHOP, BLACK_KNIGHT, BLACK_ROOK}
	for i := 0; i < 8; i++ {
		board[1][i] = WHITE_PAWN
	}
	for i := 0; i < 8; i++ {
		board[6][i] = BLACK_PAWN
	}
	return nil
}

func (ne *NoobEngine) Run() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		if !ne.turn {
			fmt.Println("Enter White's Move: ")
		} else {
			fmt.Println("Enter Black's Move: ")
		}
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
		err = ne.MakeMove(piece, init, final)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (ne *NoobEngine) PieceFromNotation(n string) (Piece, error) {
	if ne.turn {
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

func (ne *NoobEngine) IsLegal(p Piece, init Position, final Position) (bool, error) {
	if init.Equal(final) {
		return false, errors.New("init and final are same")
	}

	// check if final pos have same color piece
	if ne.board[final.GetRank()][final.GetFile()] != -1 {
		if p.Color() == ne.board[final.GetRank()][final.GetFile()].Color() {
			return false, errors.New("attacking same color")
		}
	}

	if p.IsKing() {
		if !KingAttack(&ne.board, p, init, final) {
			return false, errors.New("King cant move like that")
		}
		if ne.IsAttackedBySide(final, !p.Color()) {
			return false, errors.New("King can't move into check")
		}
	} else if p.IsPawn() {
		if !PawnMove(&ne.board, p, init, final) {
			return false, errors.New("Pawn can't move there")
		}
	} else {
		if !ne.isAttackedByPiece[p](&ne.board, p, init, final) {
			return false, errors.New("piece can't go there")
		}
	}

	if ne.inCheck {

		tarPiece := ne.board.Get(final)
		ne.board.MakeMove(p, init, final)
		var kingPos Position
		if p.Color() {
			kingPos = ne.board.GetPiecePositions(BLACK_KING)[0]
		} else {
			kingPos = ne.board.GetPiecePositions(WHITE_KING)[0]
		}
		if ne.IsAttackedBySide(kingPos, !p.Color()) {
			return false, errors.New("in check!")
		}
		ne.board.Place(tarPiece, final)
		ne.board.Place(p, init)
	}
	return true, nil
}

// checks if pos is attacked by side
func (ne *NoobEngine) IsAttackedBySide(pos Position, side bool) bool {
	startingPiece := WHITE_KING
	if side {
		startingPiece = BLACK_KING
	}
	for piece := startingPiece; piece < startingPiece+6; piece += 1 {
		if ne.IsAttackedByPiece(piece, pos) {
			return true
		}
	}
	return false
}

func (ne *NoobEngine) IsAttackedByPiece(piece Piece, tar Position) bool {
	var i, j int16
	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			if ne.board[i][j] == piece {
				if i == tar.GetRank() && j == tar.GetFile() {
					continue
				}
				var pos Position = &Pos{i, j}
				if ne.isAttackedByPiece[piece](&ne.board, piece, pos, tar) {
					return true
				}
			}
		}
	}
	return false
}

func (ne *NoobEngine) MakeMoveByNotation(move string) error {
	return nil
}

func (ne *NoobEngine) MakeMove(p Piece, init Position, final Position) error {
	if legal, err := ne.IsLegal(p, init, final); !legal {
		return errors.New("illegal move: " + err.Error())
	}
	ne.board.MakeMove(p, init, final)
	ne.turn = !ne.turn
	oppositionKing := WHITE_KING
	if !p.Color() {
		oppositionKing = BLACK_KING
	}
	pos := ne.board.GetPiecePositions(oppositionKing)[0]
	if ne.IsAttackedBySide(pos, p.Color()) {
		ne.inCheck = true
	} else {
		ne.inCheck = false
	}
	ne.PrintBoard()
	return nil
}

func (ne *NoobEngine) PrintBoard() {
	fmt.Printf("xx")
	for i := 7; i >= 0; i-- {
		fmt.Printf("xxx")
	}
	fmt.Println()
	for i := 7; i >= 0; i-- {
		fmt.Printf("X")
		for j := 0; j < 8; j++ {
			if ne.board[i][j] == -1 {
				fmt.Printf(" - ")
			} else {
				fmt.Printf(" " + pieceUnicodes[ne.board[i][j]] + " ")
			}
		}
		fmt.Printf("X")
		fmt.Printf("\n")
	}
	fmt.Printf("xx")
	for i := 7; i >= 0; i-- {
		fmt.Printf("xxx")
	}
	fmt.Println()
	if ne.inCheck {
		if ne.turn {
			fmt.Printf("Black ")
		} else {
			fmt.Printf("White ")
		}
		fmt.Printf("in check!")
	}
	fmt.Println()
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
				if board[orig.GetRank()][start] != -1 {
					return false
				}
			}
		} else {
			start := min(tar.GetRank(), orig.GetRank()) + 1
			for ; start < max(tar.GetRank(), orig.GetRank()); start++ {
				if board[start][orig.GetRank()] != -1 {
					return false
				}
			}
		}
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
			if board[rank][file] != -1 {
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
