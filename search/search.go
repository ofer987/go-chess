package search

import (
	"chess/board"
	"chess/eval"
	"chess/moves"
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

// BestMove searches to the given depth and returns the best move for the side to move.
func BestMove(b *board.Board, depth int) Result {
	legal := moves.Legal(b)
	if len(legal) == 0 {
		return Result{}
	}

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

// negamax implements the negamax variant of minimax with alpha-beta pruning.
// Returns the score from the perspective of the side to move at this node.
func negamax(b *board.Board, depth, alpha, beta int) int {
	if depth == 0 {
		score := eval.Evaluate(b)
		if b.Turn == board.Black {
			return -score
		}
		return score
	}

	legal := moves.Legal(b)
	if len(legal) == 0 {
		if moves.InCheck(b) {
			// Checkmate: prefer faster mates (higher depth remaining = fewer moves to mate).
			return -(checkmate + depth)
		}
		return 0 // stalemate
	}

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
