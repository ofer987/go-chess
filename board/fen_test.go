package board

import "testing"

func TestParseFEN(t *testing.T) {
	// Starting position: verify key piece placements and board state.
	start, err := ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if err != nil {
		t.Fatalf("ParseFEN(start): %v", err)
	}
	if got := start.Squares[0]; got != (Piece{Rook, White}) {
		t.Errorf("start a1 = %v, want White Rook", got)
	}
	if got := start.Squares[4]; got != (Piece{King, White}) {
		t.Errorf("start e1 = %v, want White King", got)
	}
	if got := start.Squares[60]; got != (Piece{King, Black}) {
		t.Errorf("start e8 = %v, want Black King", got)
	}
	if start.Turn != White {
		t.Errorf("start.Turn = %v, want White", start.Turn)
	}
	if !start.WhiteKingSide || !start.WhiteQueenSide || !start.BlackKingSide || !start.BlackQueenSide {
		t.Error("start castling rights not all true")
	}
	if start.EnPassant != -1 {
		t.Errorf("start.EnPassant = %d, want -1", start.EnPassant)
	}
	if start.HalfMove != 0 || start.FullMove != 1 {
		t.Errorf("start clocks = %d/%d, want 0/1", start.HalfMove, start.FullMove)
	}

	// After 1.e4: en passant target set, turn flipped.
	afterE4, err := ParseFEN("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1")
	if err != nil {
		t.Fatalf("ParseFEN(after e4): %v", err)
	}
	if got := afterE4.Squares[28]; got != (Piece{Pawn, White}) {
		t.Errorf("after e4: e4 = %v, want White Pawn", got)
	}
	if afterE4.Turn != Black {
		t.Errorf("after e4: Turn = %v, want Black", afterE4.Turn)
	}
	if afterE4.EnPassant != 20 { // e3
		t.Errorf("after e4: EnPassant = %d, want 20 (e3)", afterE4.EnPassant)
	}

	// No castling rights.
	noCastle, err := ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w - - 0 1")
	if err != nil {
		t.Fatalf("ParseFEN(no castle): %v", err)
	}
	if noCastle.WhiteKingSide || noCastle.WhiteQueenSide || noCastle.BlackKingSide || noCastle.BlackQueenSide {
		t.Error("no-castle FEN: expected all castling rights false")
	}

	// Move counters.
	counted, err := ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 5 10")
	if err != nil {
		t.Fatalf("ParseFEN(counters): %v", err)
	}
	if counted.HalfMove != 5 || counted.FullMove != 10 {
		t.Errorf("counters: clocks = %d/%d, want 5/10", counted.HalfMove, counted.FullMove)
	}

	// Invalid inputs.
	invalid := []string{
		"rnbqkbnr w KQkq", // too few fields
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR x KQkq - 0 1",   // bad active color
		"rnbXkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",   // unknown piece char
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq z9 0 1",  // bad en passant
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - abc 1", // non-numeric halfmove
	}

	for _, fen := range invalid {
		if _, err := ParseFEN(fen); err == nil {
			t.Errorf("ParseFEN(%q) expected error, got nil", fen)
		}
	}
}

func TestParseSquare(t *testing.T) {
	valid := []struct {
		input string
		want  int
	}{
		{"a1", 0},
		{"h1", 7},
		{"a8", 56},
		{"h8", 63},
		{"e4", 28},
		{"d5", 35},
		{"a2", 8},
		{"h7", 55},
	}

	for _, tc := range valid {
		got, err := ParseSquare(tc.input)
		if err != nil {
			t.Errorf("ParseSquare(%q) returned unexpected error: %v", tc.input, err)
			continue
		}

		if got != tc.want {
			t.Errorf("ParseSquare(%q) = %d, want %d", tc.input, got, tc.want)
		}
	}

	invalid := []string{
		"",
		"a",
		"a12",
		"i1", // file out of range
		"a9", // rank out of range
		"a0", // rank out of range
		"A1", // uppercase file
		"  ", // spaces
	}

	for _, input := range invalid {
		if _, err := ParseSquare(input); err == nil {
			t.Errorf("ParseSquare(%q) expected an error, got nil", input)
		}
	}
}

func TestSquareName(t *testing.T) {
	cases := []struct {
		sq   int
		want string
	}{
		{0, "a1"},
		{7, "h1"},
		{56, "a8"},
		{63, "h8"},
		{28, "e4"},
		{35, "d5"},
		{8, "a2"},
		{55, "h7"},
	}

	for _, tc := range cases {
		got := SquareName(tc.sq)
		if got != tc.want {
			t.Errorf("SquareName(%d) = %q, want %q", tc.sq, got, tc.want)
		}
	}
}

func TestParseSquareAndSquareNameRoundTrip(t *testing.T) {
	for sq := 0; sq < 64; sq++ {
		name := SquareName(sq)
		got, err := ParseSquare(name)
		if err != nil {
			t.Errorf("ParseSquare(SquareName(%d)) = %q, unexpected error: %v", sq, name, err)
			continue
		}

		if got != sq {
			t.Errorf("ParseSquare(SquareName(%d)) = %d, want %d", sq, got, sq)
		}
	}
}
