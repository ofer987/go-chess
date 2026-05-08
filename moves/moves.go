package moves

import "chess/board"

// Move represents a half-move (ply).
type Move struct {
	From      int
	To        int
	Promotion board.PieceType // Empty if not a promotion
}

func (m Move) String() string {
	s := board.SquareName(m.From) + board.SquareName(m.To)
	switch m.Promotion {
	case board.Queen:
		s += "q"
	case board.Rook:
		s += "r"
	case board.Bishop:
		s += "b"
	case board.Knight:
		s += "n"
	}

	return s
}

// Legal returns all legal moves for the side to move.
func Legal(b *board.Board) []Move {
	pseudo := pseudoLegal(b)
	legal := pseudo[:0]
	for _, m := range pseudo {
		nb := Apply(b, m)
		if !inCheck(nb, b.Turn) {
			legal = append(legal, m)
		}
	}

	return legal
}

// InCheck reports whether the side to move is in check.
func InCheck(b *board.Board) bool {
	return inCheck(b, b.Turn)
}

// Apply returns a new board state after making move m (no legality validation).
func Apply(b *board.Board, m Move) *board.Board {
	nb := *b
	piece := nb.Squares[m.From]
	nb.Squares[m.From] = board.NoPiece

	if m.Promotion != board.Empty {
		piece = board.Piece{Type: m.Promotion, Color: piece.Color}
	}

	// En passant: remove the captured pawn sitting one rank behind the target square.
	if piece.Type == board.Pawn && m.To == b.EnPassant {
		if piece.Color == board.White {
			nb.Squares[m.To-8] = board.NoPiece
		} else {
			nb.Squares[m.To+8] = board.NoPiece
		}
	}

	// Castling: also move the rook.
	if piece.Type == board.King {
		switch {
		case m.From == 4 && m.To == 6:
			nb.Squares[5] = nb.Squares[7]
			nb.Squares[7] = board.NoPiece
		case m.From == 4 && m.To == 2:
			nb.Squares[3] = nb.Squares[0]
			nb.Squares[0] = board.NoPiece
		case m.From == 60 && m.To == 62:
			nb.Squares[61] = nb.Squares[63]
			nb.Squares[63] = board.NoPiece
		case m.From == 60 && m.To == 58:
			nb.Squares[59] = nb.Squares[56]
			nb.Squares[56] = board.NoPiece
		}
	}

	nb.Squares[m.To] = piece

	// Update castling rights.
	if piece.Type == board.King {
		if piece.Color == board.White {
			nb.WhiteKingSide = false
			nb.WhiteQueenSide = false
		} else {
			nb.BlackKingSide = false
			nb.BlackQueenSide = false
		}
	}
	if m.From == 0 || m.To == 0 {
		nb.WhiteQueenSide = false
	}
	if m.From == 7 || m.To == 7 {
		nb.WhiteKingSide = false
	}
	if m.From == 56 || m.To == 56 {
		nb.BlackQueenSide = false
	}
	if m.From == 63 || m.To == 63 {
		nb.BlackKingSide = false
	}

	// Update en passant target.
	nb.EnPassant = -1
	if piece.Type == board.Pawn {
		if piece.Color == board.White && m.To-m.From == 16 {
			nb.EnPassant = m.From + 8
		} else if piece.Color == board.Black && m.From-m.To == 16 {
			nb.EnPassant = m.From - 8
		}
	}

	if piece.Type == board.Pawn || b.Squares[m.To].Type != board.Empty {
		nb.HalfMove = 0
	} else {
		nb.HalfMove++
	}

	if b.Turn == board.Black {
		nb.FullMove++
	}

	nb.Turn = board.Opposite(b.Turn)

	return &nb
}

func pseudoLegal(b *board.Board) []Move {
	var ms []Move
	for sq := 0; sq < 64; sq++ {
		p := b.Squares[sq]
		if p.Type == board.Empty || p.Color != b.Turn {
			continue
		}

		switch p.Type {
		case board.Pawn:
			ms = append(ms, pawnMoves(b, sq, p.Color)...)
		case board.Knight:
			ms = append(ms, knightMoves(b, sq, p.Color)...)
		case board.Bishop:
			ms = append(ms, slidingMoves(b, sq, p.Color, bishopDirs)...)
		case board.Rook:
			ms = append(ms, slidingMoves(b, sq, p.Color, rookDirs)...)
		case board.Queen:
			ms = append(ms, slidingMoves(b, sq, p.Color, queenDirs)...)
		case board.King:
			ms = append(ms, kingMoves(b, sq, p.Color)...)
		}
	}

	return ms
}

var bishopDirs = [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
var rookDirs = [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
var queenDirs = [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}

func rankOf(sq int) int { return sq / 8 }
func fileOf(sq int) int { return sq % 8 }

func slidingMoves(b *board.Board, from int, color board.Color, dirs [][2]int) []Move {
	var ms []Move
	for _, d := range dirs {
		r, f := rankOf(from)+d[0], fileOf(from)+d[1]
		for r >= 0 && r <= 7 && f >= 0 && f <= 7 {
			sq := r*8 + f
			target := b.Squares[sq]
			if target.Type == board.Empty {
				ms = append(ms, Move{from, sq, board.Empty})
			} else {
				if target.Color != color {
					ms = append(ms, Move{from, sq, board.Empty})
				}
				break
			}
			r += d[0]
			f += d[1]
		}
	}

	return ms
}

func knightMoves(b *board.Board, from int, color board.Color) []Move {
	var ms []Move
	deltas := [][2]int{{2, 1}, {2, -1}, {-2, 1}, {-2, -1}, {1, 2}, {1, -2}, {-1, 2}, {-1, -2}}
	for _, d := range deltas {
		r, f := rankOf(from)+d[0], fileOf(from)+d[1]
		if r < 0 || r > 7 || f < 0 || f > 7 {
			continue
		}

		sq := r*8 + f
		if b.Squares[sq].Color != color {
			ms = append(ms, Move{from, sq, board.Empty})
		}
	}

	return ms
}

var promotions = []board.PieceType{board.Queen, board.Rook, board.Bishop, board.Knight}

func pawnMoves(b *board.Board, from int, color board.Color) []Move {
	var ms []Move
	dir := 1
	startRank := 1
	promRank := 7
	if color == board.Black {
		dir = -1
		startRank = 6
		promRank = 0
	}

	r, f := rankOf(from), fileOf(from)

	// Single push.
	nr := r + dir
	if nr >= 0 && nr <= 7 {
		sq := nr*8 + f
		if b.Squares[sq].Type == board.Empty {
			if nr == promRank {
				for _, pt := range promotions {
					ms = append(ms, Move{from, sq, pt})
				}
			} else {
				ms = append(ms, Move{from, sq, board.Empty})
				// Double push from starting rank.
				if r == startRank {
					sq2 := (nr+dir)*8 + f
					if b.Squares[sq2].Type == board.Empty {
						ms = append(ms, Move{from, sq2, board.Empty})
					}
				}
			}
		}
	}

	// Diagonal captures (including en passant).
	for _, df := range []int{-1, 1} {
		nf := f + df
		nr := r + dir
		if nf < 0 || nf > 7 || nr < 0 || nr > 7 {
			continue
		}

		sq := nr*8 + nf
		target := b.Squares[sq]
		if (target.Type != board.Empty && target.Color != color) || sq == b.EnPassant {
			if nr == promRank {
				for _, pt := range promotions {
					ms = append(ms, Move{from, sq, pt})
				}
			} else {
				ms = append(ms, Move{from, sq, board.Empty})
			}
		}
	}

	return ms
}

func kingMoves(b *board.Board, from int, color board.Color) []Move {
	var ms []Move
	for _, d := range queenDirs {
		r, f := rankOf(from)+d[0], fileOf(from)+d[1]
		if r < 0 || r > 7 || f < 0 || f > 7 {
			continue
		}

		sq := r*8 + f
		if b.Squares[sq].Color != color {
			ms = append(ms, Move{from, sq, board.Empty})
		}
	}

	// Castling: king must not be in check, and must not pass through an attacked square.
	if color == board.White && from == 4 && !inCheck(b, color) {
		if b.WhiteKingSide &&
			b.Squares[5].Type == board.Empty &&
			b.Squares[6].Type == board.Empty &&
			!squareAttacked(b, 5, board.Black) {
			ms = append(ms, Move{4, 6, board.Empty})
		}
		if b.WhiteQueenSide &&
			b.Squares[3].Type == board.Empty &&
			b.Squares[2].Type == board.Empty &&
			b.Squares[1].Type == board.Empty &&
			!squareAttacked(b, 3, board.Black) {
			ms = append(ms, Move{4, 2, board.Empty})
		}
	}
	if color == board.Black && from == 60 && !inCheck(b, color) {
		if b.BlackKingSide &&
			b.Squares[61].Type == board.Empty &&
			b.Squares[62].Type == board.Empty &&
			!squareAttacked(b, 61, board.White) {
			ms = append(ms, Move{60, 62, board.Empty})
		}
		if b.BlackQueenSide &&
			b.Squares[59].Type == board.Empty &&
			b.Squares[58].Type == board.Empty &&
			b.Squares[57].Type == board.Empty &&
			!squareAttacked(b, 59, board.White) {
			ms = append(ms, Move{60, 58, board.Empty})
		}
	}

	return ms
}

func inCheck(b *board.Board, color board.Color) bool {
	for sq := 0; sq < 64; sq++ {
		p := b.Squares[sq]
		if p.Type == board.King && p.Color == color {
			return squareAttacked(b, sq, board.Opposite(color))
		}
	}

	return true // king missing — treat as check
}

// squareAttacked reports whether sq is attacked by any piece of the given color.
func squareAttacked(b *board.Board, sq int, by board.Color) bool {
	// Knights
	for _, d := range [][2]int{{2, 1}, {2, -1}, {-2, 1}, {-2, -1}, {1, 2}, {1, -2}, {-1, 2}, {-1, -2}} {
		r, f := rankOf(sq)+d[0], fileOf(sq)+d[1]
		if r < 0 || r > 7 || f < 0 || f > 7 {
			continue
		}

		p := b.Squares[r*8+f]
		if p.Type == board.Knight && p.Color == by {
			return true
		}
	}

	// Diagonals (bishops and queens)
	for _, d := range bishopDirs {
		r, f := rankOf(sq)+d[0], fileOf(sq)+d[1]
		for r >= 0 && r <= 7 && f >= 0 && f <= 7 {
			p := b.Squares[r*8+f]
			if p.Type != board.Empty {
				if p.Color == by && (p.Type == board.Bishop || p.Type == board.Queen) {
					return true
				}
				break
			}
			r += d[0]
			f += d[1]
		}
	}

	// Straights (rooks and queens)
	for _, d := range rookDirs {
		r, f := rankOf(sq)+d[0], fileOf(sq)+d[1]
		for r >= 0 && r <= 7 && f >= 0 && f <= 7 {
			p := b.Squares[r*8+f]
			if p.Type != board.Empty {
				if p.Color == by && (p.Type == board.Rook || p.Type == board.Queen) {
					return true
				}
				break
			}
			r += d[0]
			f += d[1]
		}
	}

	// King
	for _, d := range queenDirs {
		r, f := rankOf(sq)+d[0], fileOf(sq)+d[1]
		if r < 0 || r > 7 || f < 0 || f > 7 {
			continue
		}

		p := b.Squares[r*8+f]
		if p.Type == board.King && p.Color == by {
			return true
		}
	}

	// Pawns: a white pawn at (rank-1, file±1) attacks sq; black pawn at (rank+1, file±1).
	pawnRankDelta := -1 // look below sq for an attacking white pawn
	if by == board.Black {
		pawnRankDelta = 1
	}

	for _, df := range []int{-1, 1} {
		r, f := rankOf(sq)+pawnRankDelta, fileOf(sq)+df
		if r < 0 || r > 7 || f < 0 || f > 7 {
			continue
		}

		p := b.Squares[r*8+f]
		if p.Type == board.Pawn && p.Color == by {
			return true
		}
	}

	return false
}
