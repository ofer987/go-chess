package moves

import (
	"chess/board"
	"testing"
)

func mustParseFEN(t *testing.T, fen string) *board.Board {
	t.Helper()
	b, err := board.ParseFEN(fen)
	if err != nil {
		t.Fatalf("ParseFEN(%q): %v", fen, err)
	}

	return b
}

const startFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func TestMoveString(t *testing.T) {
	cases := []struct {
		move Move
		want string
	}{
		{Move{From: 12, To: 28}, "e2e4"},
		{Move{From: 0, To: 7}, "a1h1"},
		{Move{From: 48, To: 56, Promotion: board.Queen, Kind: Promotion}, "a7a8q"},
		{Move{From: 48, To: 56, Promotion: board.Rook, Kind: Promotion}, "a7a8r"},
		{Move{From: 48, To: 56, Promotion: board.Bishop, Kind: Promotion}, "a7a8b"},
		{Move{From: 48, To: 56, Promotion: board.Knight, Kind: Promotion}, "a7a8n"},
	}

	for _, tc := range cases {
		got := tc.move.String()
		if got != tc.want {
			t.Errorf("Move{From:%d,To:%d,Promotion:%v,Kind:%v}.String() = %q, want %q",
				tc.move.From, tc.move.To, tc.move.Promotion, tc.move.Kind, got, tc.want)
		}
	}
}

func TestLegal(t *testing.T) {
	// Starting position: exactly 20 legal moves.
	b := mustParseFEN(t, startFEN)
	if got := Legal(b); len(got) != 20 {
		t.Errorf("Legal(start) = %d moves, want 20", len(got))
	}

	// Fool's mate: White is checkmated, 0 legal moves.
	b = mustParseFEN(t, "rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 1 3")
	if got := Legal(b); len(got) != 0 {
		t.Errorf("Legal(fool's mate) = %d moves, want 0", len(got))
	}

	// Stalemate: Black has no legal moves and is not in check.
	b = mustParseFEN(t, "k7/8/1Q6/8/8/8/8/K7 b - - 0 1")
	if got := Legal(b); len(got) != 0 {
		t.Errorf("Legal(stalemate) = %d moves, want 0", len(got))
	}

	// Both White castling moves available.
	b = mustParseFEN(t, "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1")
	ms := Legal(b)
	hasKS, hasQS := false, false
	for _, m := range ms {
		if m.Kind == KingsideCastle {
			hasKS = true
		}
		if m.Kind == QueensideCastle {
			hasQS = true
		}
	}
	if !hasKS {
		t.Error("Legal(castle): missing White kingside castling e1g1")
	}
	if !hasQS {
		t.Error("Legal(castle): missing White queenside castling e1c1")
	}

	// En passant capture available: White pawn at d5 can capture on e6.
	b = mustParseFEN(t, "rnbqkbnr/pppp1ppp/8/3Pp3/8/8/PPP1PPPP/RNBQKBNR w KQkq e6 0 2")
	ms = Legal(b)
	hasEP := false
	for _, m := range ms {
		if m.Kind == EnPassant && m.From == 35 && m.To == 44 {
			hasEP = true
			break
		}
	}
	if !hasEP {
		t.Error("Legal(en passant): missing d5xe6 capture")
	}
}

func TestInCheck(t *testing.T) {
	cases := []struct {
		name string
		fen  string
		want bool
	}{
		{
			"starting position",
			startFEN,
			false,
		},
		{
			"king attacked along rank by queen",
			"K6q/8/8/8/8/8/8/7k w - - 0 1",
			true,
		},
		{
			"fool's mate: white king attacked on diagonal",
			"rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 1 3",
			true,
		},
	}

	for _, tc := range cases {
		b := mustParseFEN(t, tc.fen)
		got := InCheck(b)
		if got != tc.want {
			t.Errorf("InCheck(%s) = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestLegalCaptures(t *testing.T) {
	// Starting position: no captures available.
	b := mustParseFEN(t, startFEN)
	if got := LegalCaptures(b); len(got) != 0 {
		t.Errorf("LegalCaptures(start) = %d moves, want 0", len(got))
	}

	// After 1.e4 e5 2.d4: only White pawn d4xe5 is available.
	b = mustParseFEN(t, "rnbqkbnr/pppp1ppp/8/4p3/3PP3/8/PPP2PPP/RNBQKBNR w KQkq - 0 3")
	ms := LegalCaptures(b)
	hasCapture := false
	for _, m := range ms {
		if m.From == 27 && m.To == 36 { // d4xe5
			hasCapture = true
		}
	}
	if !hasCapture {
		t.Error("LegalCaptures: missing d4xe5 capture")
	}
	for _, m := range ms {
		if m.Kind != Capture && m.Kind != EnPassant {
			t.Errorf("LegalCaptures returned non-capture move %v (Kind=%v)", m, m.Kind)
		}
	}

	// En passant is included.
	b = mustParseFEN(t, "rnbqkbnr/pppp1ppp/8/3Pp3/8/8/PPP1PPPP/RNBQKBNR w KQkq e6 0 2")
	ms = LegalCaptures(b)
	hasEP := false
	for _, m := range ms {
		if m.From == 35 && m.To == 44 { // d5xe6 en passant
			hasEP = true
		}
	}
	if !hasEP {
		t.Error("LegalCaptures: missing en passant capture d5xe6")
	}
}

func TestApply(t *testing.T) {
	// Pawn double push sets en passant target and resets halfmove.
	b := mustParseFEN(t, startFEN)
	nb := Apply(b, Move{From: 12, To: 28, Kind: Quiet}) // e2e4
	if nb.Squares[28] != (board.Piece{Type: board.Pawn, Color: board.White}) {
		t.Errorf("Apply e2e4: e4 = %v, want White Pawn", nb.Squares[28])
	}
	if nb.Squares[12].Type != board.Empty {
		t.Error("Apply e2e4: e2 not empty")
	}
	if nb.EnPassant != 20 {
		t.Errorf("Apply e2e4: EnPassant = %d, want 20 (e3)", nb.EnPassant)
	}
	if nb.HalfMove != 0 {
		t.Errorf("Apply e2e4: HalfMove = %d, want 0", nb.HalfMove)
	}
	if nb.Turn != board.Black {
		t.Errorf("Apply e2e4: Turn = %v, want Black", nb.Turn)
	}

	// Non-pawn quiet move increments halfmove clock.
	nb2 := Apply(b, Move{From: 1, To: 16, Kind: Quiet}) // Nb1-a3
	if nb2.HalfMove != 1 {
		t.Errorf("Apply Nb1a3: HalfMove = %d, want 1", nb2.HalfMove)
	}

	// En passant capture removes the captured pawn.
	epBoard := mustParseFEN(t, "rnbqkbnr/pppp1ppp/8/3Pp3/8/8/PPP1PPPP/RNBQKBNR w KQkq e6 0 2")
	nb3 := Apply(epBoard, Move{From: 35, To: 44, Kind: EnPassant}) // d5xe6
	if nb3.Squares[44] != (board.Piece{Type: board.Pawn, Color: board.White}) {
		t.Errorf("Apply d5xe6: e6 = %v, want White Pawn", nb3.Squares[44])
	}
	if nb3.Squares[36].Type != board.Empty {
		t.Error("Apply d5xe6: e5 not empty (captured pawn not removed)")
	}

	// Kingside castling moves king and rook and clears castling rights.
	castleBoard := mustParseFEN(t, "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1")
	nb4 := Apply(castleBoard, Move{From: 4, To: 6, Kind: KingsideCastle}) // O-O
	if nb4.Squares[6] != (board.Piece{Type: board.King, Color: board.White}) {
		t.Errorf("Apply O-O: g1 = %v, want White King", nb4.Squares[6])
	}
	if nb4.Squares[5] != (board.Piece{Type: board.Rook, Color: board.White}) {
		t.Errorf("Apply O-O: f1 = %v, want White Rook", nb4.Squares[5])
	}
	if nb4.WhiteKingSide || nb4.WhiteQueenSide {
		t.Error("Apply O-O: White castling rights not cleared")
	}

	// Pawn promotion.
	promBoard := mustParseFEN(t, "8/P7/8/8/8/8/8/K6k w - - 0 1")
	nb5 := Apply(promBoard, Move{From: 48, To: 56, Promotion: board.Queen, Kind: Promotion}) // a7a8=Q
	if nb5.Squares[56] != (board.Piece{Type: board.Queen, Color: board.White}) {
		t.Errorf("Apply a7a8q: a8 = %v, want White Queen", nb5.Squares[56])
	}
	if nb5.Squares[48].Type != board.Empty {
		t.Error("Apply a7a8q: a7 not empty")
	}

	// Black move increments FullMove.
	afterE4 := mustParseFEN(t, "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1")
	nb6 := Apply(afterE4, Move{From: 52, To: 36, Kind: Quiet}) // e7e5
	if nb6.FullMove != 2 {
		t.Errorf("Apply e7e5: FullMove = %d, want 2", nb6.FullMove)
	}
}

func TestApplyEnPassant(t *testing.T) {
	// White double push sets en passant target.
	b := mustParseFEN(t, startFEN)
	nb := Apply(b, Move{From: 12, To: 28, Kind: Quiet}) // e2e4
	if nb.EnPassant != 20 {                              // e3
		t.Errorf("after e2e4: EnPassant = %d, want 20 (e3)", nb.EnPassant)
	}

	// En passant target clears after a non-double-push move.
	nb2 := Apply(nb, Move{From: 62, To: 45, Kind: Quiet}) // Ng8f6 (Black quiet move)
	if nb2.EnPassant != -1 {
		t.Errorf("after Ng8f6: EnPassant = %d, want -1", nb2.EnPassant)
	}

	// Black double push sets en passant target.
	afterE4 := mustParseFEN(t, "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1")
	nb3 := Apply(afterE4, Move{From: 52, To: 36, Kind: Quiet}) // e7e5
	if nb3.EnPassant != 44 {                                    // e6
		t.Errorf("after e7e5: EnPassant = %d, want 44 (e6)", nb3.EnPassant)
	}

	// Black en passant capture: Black pawn at d4 captures White pawn that just played e2e4.
	// Black pawn d4 (sq 27), White pawn e4 (sq 28), en passant target e3 (sq 20).
	epBlack := mustParseFEN(t, "rnbqkbnr/ppp1pppp/8/8/3pP3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 2")
	nb4 := Apply(epBlack, Move{From: 27, To: 20, Kind: EnPassant}) // d4xe3
	if nb4.Squares[20] != (board.Piece{Type: board.Pawn, Color: board.Black}) {
		t.Errorf("Apply d4xe3: e3 = %v, want Black Pawn", nb4.Squares[20])
	}
	if nb4.Squares[28].Type != board.Empty {
		t.Error("Apply d4xe3: e4 not empty (captured White pawn not removed)")
	}
	if nb4.EnPassant != -1 {
		t.Errorf("Apply d4xe3: EnPassant = %d, want -1 (cleared after capture)", nb4.EnPassant)
	}
}
