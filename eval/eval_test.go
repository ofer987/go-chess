package eval

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

func TestEvaluate(t *testing.T) {
	// Starting position is symmetric: score must be 0.
	b := mustParseFEN(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if got := Evaluate(b); got != 0 {
		t.Errorf("Evaluate(start) = %d, want 0", got)
	}

	// White up a pawn (Black's a7 pawn removed): score must be positive.
	b = mustParseFEN(t, "rnbqkbnr/1ppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if got := Evaluate(b); got <= 0 {
		t.Errorf("Evaluate(White up pawn) = %d, want > 0", got)
	}

	// Black up a queen (White's d1 queen removed): score must be negative.
	b = mustParseFEN(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNB1KBNR w - - 0 1")
	if got := Evaluate(b); got >= 0 {
		t.Errorf("Evaluate(Black up queen) = %d, want < 0", got)
	}
}
