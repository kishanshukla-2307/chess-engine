package chessEngine

import (
	"errors"
	"fmt"
)

type Board struct {
	squares [8][8]Piece
	inCheck bool
	turn    bool
}

func (b *Board) GetPiecePositions(piece Piece) []Position {
	var res []Position
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if piece == b.squares[i][j] {
				res = append(res, &Pos{(int16)(i), (int16)(j)})
			}
		}
	}
	return res
}

func (b *Board) Get(pos Position) Piece {
	return b.squares[pos.GetRank()][pos.GetFile()]
}

func (b *Board) Place(piece Piece, pos Position) {
	b.squares[pos.GetRank()][pos.GetFile()] = piece
}

func (b *Board) IsEmpty(pos Position) bool {
	return b.Get(pos) == -1
}

// func (b *Board) MakeMove(piece Piece, init Position, final Position) {
// 	b.Place(piece, final)
// 	b.Place(-1, init)
// }

func (b *Board) InitializeBoard(chess960 bool) error {
	if chess960 {
		return errors.New("chess 960 not supported yet")
	}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			b.squares[i][j] = -1
		}
	}
	b.squares[0] = [8]Piece{WHITE_ROOK, WHITE_KNIGHT, WHITE_BISHOP, WHITE_QUEEN, WHITE_KING, WHITE_BISHOP, WHITE_KNIGHT, WHITE_ROOK}
	b.squares[7] = [8]Piece{BLACK_ROOK, BLACK_KNIGHT, BLACK_BISHOP, BLACK_QUEEN, BLACK_KING, BLACK_BISHOP, BLACK_KNIGHT, BLACK_ROOK}
	for i := 0; i < 8; i++ {
		b.squares[1][i] = WHITE_PAWN
	}
	for i := 0; i < 8; i++ {
		b.squares[6][i] = BLACK_PAWN
	}
	return nil
}

func (b *Board) IsLegal(p Piece, init Position, final Position) (bool, error) {
	if b.Get(init) != p {
		return false, errors.New("No such piece there")
	}
	if init.Equal(final) {
		return false, errors.New("init and final are same")
	}

	// check if final pos have same color piece
	if b.Get(final) != -1 {
		if p.Color() == b.Get(final).Color() {
			return false, errors.New("attacking same color")
		}
	}

	if p.IsKing() {
		if !KingAttack(b, p, init, final) {
			return false, errors.New("King cant move like that")
		}
		tarPiece := b.Get(final)
		b.Place(p, final)
		b.Place(-1, init)
		if b.IsAttackedBySide(final, !p.Color()) {
			b.Place(tarPiece, final)
			b.Place(p, init)
			return false, errors.New("King can't move into check")
		}
		b.Place(tarPiece, final)
		b.Place(p, init)

	} else if p.IsPawn() {
		if !PawnMove(b, p, init, final) {
			return false, errors.New("Pawn can't move there")
		}
	} else {
		if !b.IsAttackedByPieceWithPos(p)(b, p, init, final) {
			return false, errors.New("piece can't go there")
		}
	}

	// if b.inCheck {
	tarPiece := b.Get(final)
	b.Place(p, final)
	b.Place(-1, init)
	var kingPos Position
	if p.Color() {
		kingPos = b.GetPiecePositions(BLACK_KING)[0]
	} else {
		kingPos = b.GetPiecePositions(WHITE_KING)[0]
	}
	if b.IsAttackedBySide(kingPos, !p.Color()) {
		b.Place(tarPiece, final)
		b.Place(p, init)
		return false, errors.New("in check!")
	}
	b.Place(tarPiece, final)
	b.Place(p, init)
	// }
	return true, nil
}

func (b *Board) MakeMove(p Piece, init Position, final Position) error {
	if legal, err := b.IsLegal(p, init, final); !legal {
		return errors.New("illegal move: " + err.Error())
	}
	b.Place(p, final)
	b.Place(-1, init)
	b.turn = !b.turn
	oppositionKing := WHITE_KING
	if !p.Color() {
		oppositionKing = BLACK_KING
	}
	pos := b.GetPiecePositions(oppositionKing)[0]
	if b.IsAttackedBySide(pos, p.Color()) {
		b.inCheck = true
	} else {
		b.inCheck = false
	}
	b.PrintBoard()
	return nil
}

func (b *Board) PrintBoard() {
	fmt.Printf("xx")
	for i := 7; i >= 0; i-- {
		fmt.Printf("xxx")
	}
	fmt.Println()
	for i := 7; i >= 0; i-- {
		fmt.Printf("X")
		for j := 0; j < 8; j++ {
			if b.squares[i][j] == -1 {
				fmt.Printf(" - ")
			} else {
				fmt.Printf(" " + pieceUnicodes[b.squares[i][j]] + " ")
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
	if b.inCheck {
		if b.turn {
			fmt.Printf("Black ")
		} else {
			fmt.Printf("White ")
		}
		fmt.Printf("in check!")
	}
	fmt.Println()
}

// checks if pos is attacked by side
func (b *Board) IsAttackedBySide(pos Position, side bool) bool {
	startingPiece := WHITE_KING
	if side {
		startingPiece = BLACK_KING
	}
	for piece := startingPiece; piece < startingPiece+6; piece += 1 {
		if b.IsAttackedByPiece(piece, pos) {
			return true
		}
	}
	return false
}

func (b *Board) IsAttackedByPiece(piece Piece, tar Position) bool {
	var i, j int16
	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			if b.squares[i][j] == piece {
				if i == tar.GetRank() && j == tar.GetFile() {
					continue
				}
				var pos Position = &Pos{i, j}
				if b.IsAttackedByPieceWithPos(piece)(b, piece, pos, tar) {
					return true
				}
			}
		}
	}
	return false
}

func (b *Board) IsAttackedByPieceWithPos(piece Piece) func(*Board, Piece, Position, Position) bool {
	switch {
	case piece == WHITE_KING || piece == BLACK_KING:
		return KingAttack
	case piece == WHITE_QUEEN || piece == BLACK_QUEEN:
		return QueenAttack
	case piece == WHITE_ROOK || piece == BLACK_ROOK:
		return RookAttack
	case piece == WHITE_BISHOP || piece == BLACK_BISHOP:
		return BishopAttack
	case piece == WHITE_KNIGHT || piece == BLACK_KNIGHT:
		return KnightAttack
	case piece == WHITE_PAWN || piece == BLACK_PAWN:
		return PawnAttack
	}
	return nil
}

func (b *Board) GenerateMovesForPiece(piece Piece) func(*Board, Position) []Move {
	switch {
	case piece == WHITE_KING || piece == BLACK_KING:
		return GenerateKingMoves
	case piece == WHITE_QUEEN || piece == BLACK_QUEEN:
		return GenerateQueenMoves
	case piece == WHITE_ROOK || piece == BLACK_ROOK:
		return GenerateRookMoves
	case piece == WHITE_BISHOP || piece == BLACK_BISHOP:
		return GenerateBishopMoves
	case piece == WHITE_KNIGHT || piece == BLACK_KNIGHT:
		return GenerateKnightMoves
	case piece == WHITE_PAWN || piece == BLACK_PAWN:
		return GeneratePawnMoves
	}
	return nil
}

func (b *Board) GenerateMoves(side bool) []Move {
	var moves []Move
	if side {
		for _, piece := range BLACK_PIECES {
			positions := b.GetPiecePositions(piece)
			for _, pos := range positions {
				moves = append(moves, b.GenerateMovesForPiece(piece)(b, pos)...)
			}
		}
	} else {
		for _, piece := range WHITE_PIECES {
			positions := b.GetPiecePositions(piece)
			for _, pos := range positions {
				moves = append(moves, b.GenerateMovesForPiece(piece)(b, pos)...)
			}
		}
	}
	var legalMoves []Move
	for _, move := range moves {
		if legal, err := b.IsLegal(move.piece, move.init, move.final); (err == nil) && legal {
			legalMoves = append(legalMoves, move)
		}
	}

	return legalMoves
}

func GenerateKingMoves(board *Board, pos Position) []Move {
	x := pos.GetRank()
	y := pos.GetFile()
	finalPositions := []Position{&Pos{x - 1, y}, &Pos{x + 1, y}, &Pos{x, y - 1}, &Pos{x, y + 1}, &Pos{x - 1, y - 1}, &Pos{x - 1, y + 1}, &Pos{x + 1, y + 1}, &Pos{x + 1, y - 1}}
	var moves []Move
	for _, finalPos := range finalPositions {
		if finalPos.IsValid() {
			moves = append(moves, Move{board.Get(pos), pos, finalPos})
		}
	}
	return moves
}

func GenerateRookMoves(board *Board, pos Position) []Move {
	piece := board.Get(pos)
	x := pos.GetRank()
	y := pos.GetFile()
	var finalPositions []Position
	add := func(x int16, y int16) {
		if board.Get(&Pos{x, y}) != -1 {
			if board.Get(&Pos{x, y}).Color() != piece.Color() {
				finalPositions = append(finalPositions, &Pos{x, y})
			}
			return
		}
		finalPositions = append(finalPositions, &Pos{x, y})
	}
	var itr int16
	for itr = x + 1; itr < 8; itr++ {
		add(itr, y)
	}
	for itr = y + 1; itr < 8; itr++ {
		add(x, itr)
	}
	for itr = x - 1; itr >= 0; itr-- {
		add(itr, y)
	}
	for itr = y - 1; itr >= 0; itr-- {
		add(x, itr)
	}
	var moves []Move
	for _, finalPos := range finalPositions {
		moves = append(moves, Move{board.Get(pos), pos, finalPos})
	}
	return moves
}

func GenerateBishopMoves(board *Board, pos Position) []Move {
	piece := board.Get(pos)
	x := pos.GetRank()
	y := pos.GetFile()
	var finalPositions []Position
	add := func(x int16, y int16) {
		if board.Get(&Pos{x, y}) != -1 {
			if board.Get(&Pos{x, y}).Color() != piece.Color() {
				finalPositions = append(finalPositions, &Pos{x, y})
			}
			return
		}
		finalPositions = append(finalPositions, &Pos{x, y})
	}
	var xitr, yitr int16 = x + 1, y + 1
	for xitr < 8 && yitr < 8 {
		add(xitr, yitr)
		xitr++
		yitr++
	}
	xitr, yitr = x-1, y-1
	for xitr >= 0 && yitr >= 0 {
		add(xitr, yitr)
		xitr--
		yitr--
	}
	xitr, yitr = x+1, y-1
	for xitr < 8 && yitr >= 0 {
		add(xitr, yitr)
		xitr++
		yitr--
	}
	xitr, yitr = x-1, y+1
	for xitr >= 0 && yitr < 8 {
		add(xitr, yitr)
		xitr--
		yitr++
	}
	var moves []Move
	for _, finalPos := range finalPositions {
		moves = append(moves, Move{board.Get(pos), pos, finalPos})
	}
	return moves
}

func GenerateQueenMoves(board *Board, pos Position) []Move {
	var moves []Move
	moves = append(moves, GenerateRookMoves(board, pos)...)
	moves = append(moves, GenerateBishopMoves(board, pos)...)
	return moves
}

func GenerateKnightMoves(board *Board, pos Position) []Move {
	x := pos.GetRank()
	y := pos.GetFile()
	finalPositions := []Position{&Pos{x + 1, y + 2}, &Pos{x + 1, y - 2}, &Pos{x + 2, y - 1}, &Pos{x + 2, y + 1}, &Pos{x - 1, y - 2}, &Pos{x - 1, y + 2}, &Pos{x - 2, y - 1}, &Pos{x - 2, y + 1}}
	var moves []Move
	for _, finalPos := range finalPositions {
		if finalPos.IsValid() {
			moves = append(moves, Move{board.Get(pos), pos, finalPos})
		}
	}
	return moves
}

func GeneratePawnMoves(board *Board, pos Position) []Move {
	x := pos.GetRank()
	y := pos.GetFile()
	var finalPositions []Position
	if board.Get(pos).Color() {
		pos := &Pos{x - 1, y}
		if pos.IsValid() && board.Get(pos) == -1 {
			finalPositions = append(finalPositions, &Pos{x - 1, y})
		}
		pos1 := &Pos{x - 2, y}
		if x == 6 && pos.IsValid() && board.Get(pos) == -1 && pos1.IsValid() && board.Get(pos1) == -1 {
			finalPositions = append(finalPositions, &Pos{x - 2, y})
		}
		pos2 := &Pos{x - 1, y - 1}
		if (pos2).IsValid() && !board.IsEmpty(pos2) && !board.Get(pos2).Color() {
			finalPositions = append(finalPositions, &Pos{x - 1, y - 1})
		}
		pos3 := &Pos{x - 1, y + 1}
		if (pos3).IsValid() && !board.IsEmpty(pos3) && !board.Get(pos3).Color() {
			finalPositions = append(finalPositions, &Pos{x - 1, y + 1})
		}
	} else {
		pos := &Pos{x + 1, y}
		if pos.IsValid() && board.Get(pos) == -1 {
			finalPositions = append(finalPositions, &Pos{x + 1, y})
		}
		pos1 := &Pos{x + 2, y}
		if x == 1 && pos.IsValid() && board.Get(&Pos{x + 1, y}) == -1 && pos1.IsValid() && board.Get(&Pos{x + 2, y}) == -1 {
			finalPositions = append(finalPositions, &Pos{x + 2, y})
		}
		pos2 := &Pos{x + 1, y - 1}
		if (pos2).IsValid() && !board.IsEmpty(pos2) && board.Get(pos2).Color() {
			finalPositions = append(finalPositions, &Pos{x + 1, y - 1})
		}
		pos3 := &Pos{x + 1, y + 1}
		if (pos3).IsValid() && !board.IsEmpty(pos3) && board.Get(pos3).Color() {
			finalPositions = append(finalPositions, &Pos{x + 1, y + 1})
		}
	}
	var moves []Move
	for _, finalPos := range finalPositions {
		if finalPos.IsValid() {
			moves = append(moves, Move{board.Get(pos), pos, finalPos})
		}
	}
	return moves
}

type Move struct {
	piece Piece
	init  Position
	final Position
}

type Position interface {
	GetRank() int16
	GetFile() int16
	Equal(Position) bool
	Sub(Position) Position
	IsValid() bool
}

type Pos struct {
	X int16
	Y int16
}

func (p *Pos) GetRank() int16 {
	return p.X
}

func (p *Pos) GetFile() int16 {
	return p.Y
}

func (p *Pos) Sub(pos Position) Position {
	return &Pos{p.GetRank() - pos.GetRank(), p.GetFile() - pos.GetFile()}
}

func (p *Pos) Equal(pos Position) bool {
	return p.X == pos.GetRank() && p.Y == pos.GetFile()
}

func (p *Pos) IsValid() bool {
	return (p.GetRank() <= 7 && p.GetRank() >= 0) && (p.GetFile() <= 7 && p.GetFile() >= 0)
}

var _ Position = (*Pos)(nil)
