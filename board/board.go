package board

type PieceType int8

const (
	Empty  PieceType = 0
	Pawn   PieceType = 1
	Knight PieceType = 2
	Bishop PieceType = 3
	Rook   PieceType = 4
	Queen  PieceType = 5
	King   PieceType = 6
)

type Color int8

const (
	NoColor Color = 0
	White   Color = 1
	Black   Color = 2
)

func Opposite(c Color) Color {
	if c == White {
		return Black
	}
	return White
}

type Piece struct {
	Type  PieceType
	Color Color
}

var NoPiece = Piece{Empty, NoColor}

// Board square index: rank*8 + file, a1=0, h1=7, a8=56, h8=63.
type Board struct {
	Squares        [64]Piece
	Turn           Color
	WhiteKingSide  bool
	WhiteQueenSide bool
	BlackKingSide  bool
	BlackQueenSide bool
	EnPassant      int // -1 if none
	HalfMove       int
	FullMove       int
}

func NewBoard() *Board {
	return &Board{EnPassant: -1}
}

func (b *Board) IsEmpty(sq int) bool {
	return b.Squares[sq].Type == Empty
}

func (b *Board) PieceAt(sq int) Piece {
	return b.Squares[sq]
}
