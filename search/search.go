package search

import (
	"chess/board"
	"chess/eval"
	"chess/moves"
	"sort"
)

const (
	inf       = 1_000_000_000
	checkmate = 10_000_000
)

// Result holds the best move found and its score (from the moving side's perspective).
type Result struct {
	Move  moves.Move
	Score int
}

// mvvLVATable mirrors eval.pieceValue for move ordering (no eval import cycle needed).
var mvvLVATable = [7]int{
	0,      // Empty
	100,    // Pawn
	320,    // Knight
	330,    // Bishop
	500,    // Rook
	900,    // Queen
	20_000, // King
}

// mvvLVAScore returns a move-ordering score. Captures score positively (victim*10
// minus attacker), ensuring all captures outrank quiet moves (score 0).
// En passant is treated as PxP.
func mvvLVAScore(b *board.Board, m moves.Move) int {
	victim := b.Squares[m.To].Type
	if victim == board.Empty {
		if m.To != b.EnPassant {
			return 0 // quiet move
		}
		victim = board.Pawn // en passant captures a pawn
	}
	attacker := b.Squares[m.From].Type

	return mvvLVATable[victim]*10 - mvvLVATable[attacker]
}

// orderMoves sorts ms in-place: captures first (highest MVV-LVA score first),
// quiet moves last.
func orderMoves(b *board.Board, ms []moves.Move) {
	sort.Slice(ms, func(i, j int) bool {
		return mvvLVAScore(b, ms[i]) > mvvLVAScore(b, ms[j])
	})
}

// BestMove searches to the given depth and returns the best move for the side to move.
func BestMove(b *board.Board, depth int) Result {
	legal := moves.Legal(b)
	if len(legal) == 0 {
		return Result{}
	}

	orderMoves(b, legal)

	best := Result{Score: -inf}
	alpha, beta := -inf, inf

	for _, m := range legal {
		nb := moves.Apply(b, m)
		score := -negamax(nb, depth-1, -beta, -alpha)
		if score > best.Score {
			best = Result{Move: m, Score: score}
		}
		if score > alpha {
			alpha = score
		}
	}

	return best
}

// quiesce extends the search at depth=0 by searching captures until the position
// is quiet, preventing the horizon effect on tactical sequences.
// Returns the score from the perspective of the side to move.
func quiesce(b *board.Board, alpha, beta int) int {
	standPat := eval.Evaluate(b)
	if b.Turn == board.Black {
		standPat = -standPat
	}

	if standPat >= beta {
		return beta
	}
	if standPat > alpha {
		alpha = standPat
	}

	captures := moves.LegalCaptures(b)
	orderMoves(b, captures)

	for _, m := range captures {
		nb := moves.Apply(b, m)
		score := -quiesce(nb, -beta, -alpha)
		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}

	return alpha
}

// negamax implements the negamax variant of minimax with alpha-beta pruning.
// Returns the score from the perspective of the side to move at this node.
func negamax(b *board.Board, depth, alpha, beta int) int {
	if depth == 0 {
		return quiesce(b, alpha, beta)
	}

	legal := moves.Legal(b)
	if len(legal) == 0 {
		if moves.InCheck(b) {
			// Checkmate: prefer faster mates (higher depth remaining = fewer moves to mate).
			return -(checkmate + depth)
		}

		return 0 // stalemate
	}

	orderMoves(b, legal)

	for _, m := range legal {
		nb := moves.Apply(b, m)
		score := -negamax(nb, depth-1, -beta, -alpha)
		if score >= beta {
			return beta // fail-hard cutoff
		}

		if score > alpha {
			alpha = score
		}
	}

	return alpha
}
