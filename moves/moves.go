package moves

import "chess/board"

// MoveKind classifies a Move for Apply dispatch and caller inspection.
type MoveKind int

const (
	Quiet MoveKind = iota
	Capture
	EnPassant
	KingsideCastle
	QueensideCastle
	Promotion // pawn reaches back rank — quiet push or diagonal capture both
)

// Move represents a half-move (ply).
type Move struct {
	From      board.Square
	To        board.Square
	Promotion board.PieceType // non-Empty only when Kind == Promotion
	Kind      MoveKind
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

// LegalCaptures returns all legal capturing moves (including en passant) for the side to move.
func LegalCaptures(b *board.Board) []Move {
	var ms []Move
	for _, m := range pseudoLegal(b) {
		if m.Kind != Capture && m.Kind != EnPassant {
			continue
		}
		nb := Apply(b, m)
		if !inCheck(nb, b.Turn) {
			ms = append(ms, m)
		}
	}

	return ms
}

// mover applies the square mutations for a specific MoveKind.
type mover interface {
	apply(nb *board.Board, piece board.Piece, m Move)
}

type quietMover struct{}

func (quietMover) apply(nb *board.Board, piece board.Piece, m Move) {
	nb.Squares[m.To] = piece
}

type enPassantMover struct{}

func (enPassantMover) apply(nb *board.Board, piece board.Piece, m Move) {
	switch piece.Color {
	case board.White:
		nb.Squares[m.To-8] = board.NoPiece
	case board.Black:
		nb.Squares[m.To+8] = board.NoPiece
	}
	nb.Squares[m.To] = piece
}

type kingsideCastleMover struct{}

func (kingsideCastleMover) apply(nb *board.Board, piece board.Piece, m Move) {
	switch piece.Color {
	case board.White:
		nb.Squares[5] = nb.Squares[7]
		nb.Squares[7] = board.NoPiece
	case board.Black:
		nb.Squares[61] = nb.Squares[63]
		nb.Squares[63] = board.NoPiece
	}
	nb.Squares[m.To] = piece
}

type queensideCastleMover struct{}

func (queensideCastleMover) apply(nb *board.Board, piece board.Piece, m Move) {
	switch piece.Color {
	case board.White:
		nb.Squares[3] = nb.Squares[0]
		nb.Squares[0] = board.NoPiece
	case board.Black:
		nb.Squares[59] = nb.Squares[56]
		nb.Squares[56] = board.NoPiece
	}
	nb.Squares[m.To] = piece
}

type promotionMover struct{}

func (promotionMover) apply(nb *board.Board, piece board.Piece, m Move) {
	nb.Squares[m.To] = board.Piece{Type: m.Promotion, Color: piece.Color}
}

func moverFor(kind MoveKind) mover {
	switch kind {
	case EnPassant:
		return enPassantMover{}
	case KingsideCastle:
		return kingsideCastleMover{}
	case QueensideCastle:
		return queensideCastleMover{}
	case Promotion:
		return promotionMover{}
	default:
		return quietMover{} // Quiet and Capture have identical square mutations
	}
}

// Apply returns a new board state after making move m (no legality validation).
func Apply(b *board.Board, m Move) *board.Board {
	nb := *b
	piece := nb.Squares[m.From]
	nb.Squares[m.From] = board.NoPiece

	moverFor(m.Kind).apply(&nb, piece, m)

	updateCastlingRights(&nb, m, piece)
	updateEnPassantTarget(&nb, m, piece)
	updateClocks(&nb, b, m, piece)
	nb.Turn = board.Opposite(b.Turn)

	return &nb
}

func updateCastlingRights(nb *board.Board, m Move, piece board.Piece) {
	if piece.Type == board.King {
		switch piece.Color {
		case board.White:
			nb.WhiteKingSide = false
			nb.WhiteQueenSide = false
		case board.Black:
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
}

func updateEnPassantTarget(nb *board.Board, m Move, piece board.Piece) {
	nb.EnPassant = board.NoSquare
	if piece.Type == board.Pawn {
		switch piece.Color {
		case board.White:
			if m.To-m.From == 16 {
				nb.EnPassant = m.From + 8
			}
		case board.Black:
			if m.From-m.To == 16 {
				nb.EnPassant = m.From - 8
			}
		}
	}
}

func updateClocks(nb *board.Board, b *board.Board, m Move, piece board.Piece) {
	if piece.Type == board.Pawn || m.Kind == Capture || m.Kind == EnPassant {
		nb.HalfMove = 0
	} else {
		nb.HalfMove += 1
	}

	if b.Turn == board.Black {
		nb.FullMove += 1
	}
}

func pseudoLegal(b *board.Board) []Move {
	var ms []Move
	for sq := board.Square(0); sq < 64; sq += 1 {
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

func slidingMoves(b *board.Board, from board.Square, color board.Color, dirs [][2]int) []Move {
	var ms []Move
	for _, d := range dirs {
		r, f := board.RankOf(from)+d[0], board.FileOf(from)+d[1]
		for r >= 0 && r <= 7 && f >= 0 && f <= 7 {
			sq := board.SquareOf(r, f)
			target := b.Squares[sq]
			if target.Type == board.Empty {
				ms = append(ms, Move{From: from, To: sq, Kind: Quiet})
			} else {
				if target.Color != color {
					ms = append(ms, Move{From: from, To: sq, Kind: Capture})
				}
				break
			}
			r += d[0]
			f += d[1]
		}
	}

	return ms
}

func knightMoves(b *board.Board, from board.Square, color board.Color) []Move {
	var ms []Move
	deltas := [][2]int{{2, 1}, {2, -1}, {-2, 1}, {-2, -1}, {1, 2}, {1, -2}, {-1, 2}, {-1, -2}}
	for _, d := range deltas {
		r, f := board.RankOf(from)+d[0], board.FileOf(from)+d[1]
		if r < 0 || r > 7 || f < 0 || f > 7 {
			continue
		}

		sq := board.SquareOf(r, f)
		if b.Squares[sq].Color != color {
			kind := Quiet
			if b.Squares[sq].Type != board.Empty {
				kind = Capture
			}
			ms = append(ms, Move{From: from, To: sq, Kind: kind})
		}
	}

	return ms
}

var promotionPieces = []board.PieceType{board.Queen, board.Rook, board.Bishop, board.Knight}

func pawnMoves(b *board.Board, from board.Square, color board.Color) []Move {
	var ms []Move
	var dir, startRank, promRank int
	switch color {
	case board.White:
		dir, startRank, promRank = 1, 1, 7
	case board.Black:
		dir, startRank, promRank = -1, 6, 0
	}

	r, f := board.RankOf(from), board.FileOf(from)

	// Single push.
	nr := r + dir
	if nr >= 0 && nr <= 7 {
		sq := board.SquareOf(nr, f)
		if b.Squares[sq].Type == board.Empty {
			if nr == promRank {
				for _, pt := range promotionPieces {
					ms = append(ms, Move{From: from, To: sq, Promotion: pt, Kind: Promotion})
				}
			} else {
				ms = append(ms, Move{From: from, To: sq, Kind: Quiet})
				// Double push from starting rank.
				if r == startRank {
					sq2 := board.SquareOf(nr+dir, f)
					if b.Squares[sq2].Type == board.Empty {
						ms = append(ms, Move{From: from, To: sq2, Kind: Quiet})
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

		sq := board.SquareOf(nr, nf)
		target := b.Squares[sq]
		if (target.Type != board.Empty && target.Color != color) || sq == b.EnPassant {
			if nr == promRank {
				for _, pt := range promotionPieces {
					ms = append(ms, Move{From: from, To: sq, Promotion: pt, Kind: Promotion})
				}
			} else if sq == b.EnPassant {
				ms = append(ms, Move{From: from, To: sq, Kind: EnPassant})
			} else {
				ms = append(ms, Move{From: from, To: sq, Kind: Capture})
			}
		}
	}

	return ms
}

func kingMoves(b *board.Board, from board.Square, color board.Color) []Move {
	var ms []Move
	for _, d := range queenDirs {
		r, f := board.RankOf(from)+d[0], board.FileOf(from)+d[1]
		if r < 0 || r > 7 || f < 0 || f > 7 {
			continue
		}

		sq := board.SquareOf(r, f)
		if b.Squares[sq].Color != color {
			kind := Quiet
			if b.Squares[sq].Type != board.Empty {
				kind = Capture
			}
			ms = append(ms, Move{From: from, To: sq, Kind: kind})
		}
	}

	// Castling: king must not be in check, and must not pass through an attacked square.
	if color == board.White && from == 4 && !inCheck(b, color) {
		if b.WhiteKingSide &&
			b.Squares[5].Type == board.Empty &&
			b.Squares[6].Type == board.Empty &&
			!squareAttacked(b, 5, board.Black) {
			ms = append(ms, Move{From: 4, To: 6, Kind: KingsideCastle})
		}
		if b.WhiteQueenSide &&
			b.Squares[3].Type == board.Empty &&
			b.Squares[2].Type == board.Empty &&
			b.Squares[1].Type == board.Empty &&
			!squareAttacked(b, 3, board.Black) {
			ms = append(ms, Move{From: 4, To: 2, Kind: QueensideCastle})
		}
	}
	if color == board.Black && from == 60 && !inCheck(b, color) {
		if b.BlackKingSide &&
			b.Squares[61].Type == board.Empty &&
			b.Squares[62].Type == board.Empty &&
			!squareAttacked(b, 61, board.White) {
			ms = append(ms, Move{From: 60, To: 62, Kind: KingsideCastle})
		}
		if b.BlackQueenSide &&
			b.Squares[59].Type == board.Empty &&
			b.Squares[58].Type == board.Empty &&
			b.Squares[57].Type == board.Empty &&
			!squareAttacked(b, 59, board.White) {
			ms = append(ms, Move{From: 60, To: 58, Kind: QueensideCastle})
		}
	}

	return ms
}

func inCheck(b *board.Board, color board.Color) bool {
	for sq := board.Square(0); sq < 64; sq += 1 {
		p := b.Squares[sq]
		if p.Type == board.King && p.Color == color {
			return squareAttacked(b, sq, board.Opposite(color))
		}
	}

	return true // king missing — treat as check
}

// squareAttacked reports whether sq is attacked by any piece of the given color.
func squareAttacked(b *board.Board, sq board.Square, by board.Color) bool {
	// Knights
	for _, d := range [][2]int{{2, 1}, {2, -1}, {-2, 1}, {-2, -1}, {1, 2}, {1, -2}, {-1, 2}, {-1, -2}} {
		r, f := board.RankOf(sq)+d[0], board.FileOf(sq)+d[1]
		if r < 0 || r > 7 || f < 0 || f > 7 {
			continue
		}

		p := b.Squares[board.SquareOf(r, f)]
		if p.Type == board.Knight && p.Color == by {
			return true
		}
	}

	// Diagonals (bishops and queens)
	for _, d := range bishopDirs {
		r, f := board.RankOf(sq)+d[0], board.FileOf(sq)+d[1]
		for r >= 0 && r <= 7 && f >= 0 && f <= 7 {
			p := b.Squares[board.SquareOf(r, f)]
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
		r, f := board.RankOf(sq)+d[0], board.FileOf(sq)+d[1]
		for r >= 0 && r <= 7 && f >= 0 && f <= 7 {
			p := b.Squares[board.SquareOf(r, f)]
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
		r, f := board.RankOf(sq)+d[0], board.FileOf(sq)+d[1]
		if r < 0 || r > 7 || f < 0 || f > 7 {
			continue
		}

		p := b.Squares[board.SquareOf(r, f)]
		if p.Type == board.King && p.Color == by {
			return true
		}
	}

	// Pawns: a white pawn at (rank-1, file±1) attacks sq; black pawn at (rank+1, file±1).
	var pawnRankDelta int
	switch by {
	case board.White:
		pawnRankDelta = -1
	case board.Black:
		pawnRankDelta = 1
	}

	for _, df := range []int{-1, 1} {
		r, f := board.RankOf(sq)+pawnRankDelta, board.FileOf(sq)+df
		if r < 0 || r > 7 || f < 0 || f > 7 {
			continue
		}

		p := b.Squares[board.SquareOf(r, f)]
		if p.Type == board.Pawn && p.Color == by {
			return true
		}
	}

	return false
}
