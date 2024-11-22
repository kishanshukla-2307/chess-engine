package chessEngine

// interfaces

type Board [8][8]Piece

func (b *Board) GetPiecePositions(piece Piece) []Position {
	var res []Position
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if piece == b[i][j] {
				res = append(res, &Pos{(int16)(i), (int16)(j)})
			}
		}
	}
	return res
}

func (b *Board) Get(pos Position) Piece {
	return b[pos.GetRank()][pos.GetFile()]
}

func (b *Board) Place(piece Piece, pos Position) {
	b[pos.GetRank()][pos.GetFile()] = piece
}

func (b *Board) IsEmpty(pos Position) bool {
	return b.Get(pos) == -1
}

func (b *Board) MakeMove(piece Piece, init Position, final Position) {
	b.Place(piece, final)
	b.Place(-1, init)
}

type Position interface {
	GetRank() int16
	GetFile() int16
	Equal(Position) bool
	Sub(Position) Position
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

var _ Position = (*Pos)(nil)
