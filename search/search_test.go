package search

import (
	"chess/board"
	"chess/moves"
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

func TestBestMove(t *testing.T) {
	// Mate in 1: Rf7-f8# is the only checkmate. Score must reach checkmate threshold.
	// Position: White Rook f7, White King g6, Black King h8.
	b := mustParseFEN(t, "7k/5R2/6K1/8/8/8/8/8 w - - 0 1")
	result := BestMove(b, 2)
	wantMove := moves.Move{From: 53, To: 61}
	if result.Move != wantMove {
		t.Errorf("BestMove(mate-in-1) move = %v, want %v", result.Move, wantMove)
	}
	if result.Score < checkmate {
		t.Errorf("BestMove(mate-in-1) score = %d, want >= %d", result.Score, checkmate)
	}

	// Clear best capture: White queen at d4 can take free Black queen at d5.
	b = mustParseFEN(t, "4k3/8/8/3q4/3Q4/8/8/4K3 w - - 0 1")
	result = BestMove(b, 1)
	if result.Move.To != 35 { // d5
		t.Errorf("BestMove(best capture) move = %v, want To=35 (d5)", result.Move)
	}

	// No legal moves (checkmate position): returns zero Result.
	b = mustParseFEN(t, "rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 1 3")
	result = BestMove(b, 1)
	if result != (Result{}) {
		t.Errorf("BestMove(checkmate) = %v, want zero Result", result)
	}

	// Starting position at depth 1: must return a legal move (non-zero).
	b = mustParseFEN(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	result = BestMove(b, 1)
	if result.Move == (moves.Move{}) {
		t.Error("BestMove(start, depth=1) returned zero move")
	}
}
