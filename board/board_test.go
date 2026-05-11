package board

import "testing"

func TestMirrorVertical(t *testing.T) {
	cases := []struct {
		sq   Square
		want Square
	}{
		{0, 56},  // a1 → a8
		{56, 0},  // a8 → a1
		{7, 63},  // h1 → h8
		{63, 7},  // h8 → h1
		{4, 60},  // e1 → e8
		{28, 36}, // e4 → e5
	}

	for _, tc := range cases {
		got := MirrorVertical(tc.sq)
		if got != tc.want {
			t.Errorf("MirrorVertical(%d) = %d, want %d", tc.sq, got, tc.want)
		}
	}

	// Round-trip: mirroring twice returns the original square.
	for sq := Square(0); sq < 64; sq += 1 {
		if got := MirrorVertical(MirrorVertical(sq)); got != sq {
			t.Errorf("MirrorVertical(MirrorVertical(%d)) = %d, want %d", sq, got, sq)
		}
	}
}

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
