package board

import "testing"

func TestOpposite(t *testing.T) {
	cases := []struct {
		input Color
		want  Color
	}{
		{White, Black},
		{Black, White},
		{NoColor, NoColor},
	}

	for _, tc := range cases {
		got := Opposite(tc.input)
		if got != tc.want {
			t.Errorf("Opposite(%v) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestNewBoard(t *testing.T) {
	b := NewBoard()

	if b.EnPassant != NoSquare {
		t.Errorf("NewBoard().EnPassant = %d, want NoSquare", b.EnPassant)
	}

	for sq := 0; sq < 64; sq += 1 {
		if b.Squares[sq] != NoPiece {
			t.Errorf("NewBoard().Squares[%d] = %v, want NoPiece", sq, b.Squares[sq])
		}
	}
}

func TestIsEmpty(t *testing.T) {
	b := NewBoard()

	if !b.IsEmpty(0) {
		t.Error("IsEmpty(0) = false on empty board, want true")
	}

	b.Squares[0] = Piece{Rook, White}

	if b.IsEmpty(0) {
		t.Error("IsEmpty(0) = true after placing piece, want false")
	}
}

func TestPieceAt(t *testing.T) {
	b := NewBoard()
	cases := []struct {
		sq    Square
		piece Piece
	}{
		{0, Piece{Rook, White}},
		{63, Piece{King, Black}},
		{28, Piece{Pawn, White}},
	}

	for _, tc := range cases {
		b.Squares[tc.sq] = tc.piece
		got := b.PieceAt(tc.sq)
		if got != tc.piece {
			t.Errorf("PieceAt(%d) = %v, want %v", tc.sq, got, tc.piece)
		}
	}
}
