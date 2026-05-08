package board

import (
	"strings"
	"testing"
)

func TestDisplay(t *testing.T) {
	b, err := ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if err != nil {
		t.Fatalf("ParseFEN: %v", err)
	}

	out := Display(b)

	for _, want := range []string{
		"  a b c d e f g h",
		"♖", "♜",
		"White to move",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Display() missing %q", want)
		}
	}

	for _, rank := range []string{"1", "2", "3", "4", "5", "6", "7", "8"} {
		if !strings.Contains(out, rank) {
			t.Errorf("Display() missing rank label %q", rank)
		}
	}
}

func TestDisplayBlackToMove(t *testing.T) {
	b, err := ParseFEN("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1")
	if err != nil {
		t.Fatalf("ParseFEN: %v", err)
	}

	out := Display(b)

	if !strings.Contains(out, "Black to move") {
		t.Errorf("Display() does not contain \"Black to move\"\n%s", out)
	}
}
