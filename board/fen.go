package board

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseFEN parses a FEN string into a Board.
func ParseFEN(fen string) (*Board, error) {
	parts := strings.Fields(fen)
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid FEN: too few fields")
	}

	b := NewBoard()

	// Piece placement: FEN starts at rank 8 and goes down.
	rank := 7
	file := 0
	for _, ch := range parts[0] {
		if ch == '/' {
			rank -= 1
			file = 0
			continue
		}

		if ch >= '1' && ch <= '8' {
			file += int(ch - '0')
			continue
		}

		sq := SquareOf(rank, file)
		p, err := parsePieceChar(ch)
		if err != nil {
			return nil, err
		}

		b.Squares[sq] = p
		file += 1
	}

	switch parts[1] {
	case "w":
		b.Turn = White
	case "b":
		b.Turn = Black
	default:
		return nil, fmt.Errorf("invalid active color: %s", parts[1])
	}

	if parts[2] != "-" {
		for _, ch := range parts[2] {
			switch ch {
			case 'K':
				b.WhiteKingSide = true
			case 'Q':
				b.WhiteQueenSide = true
			case 'k':
				b.BlackKingSide = true
			case 'q':
				b.BlackQueenSide = true
			}
		}
	}

	b.EnPassant = -1
	if parts[3] != "-" {
		sq, err := ParseSquare(parts[3])
		if err != nil {
			return nil, err
		}

		b.EnPassant = sq
	}

	if len(parts) > 4 {
		n, err := strconv.Atoi(parts[4])
		if err != nil {
			return nil, fmt.Errorf("invalid halfmove clock: %s", parts[4])
		}

		b.HalfMove = n
	}

	if len(parts) > 5 {
		n, err := strconv.Atoi(parts[5])
		if err != nil {
			return nil, fmt.Errorf("invalid fullmove number: %s", parts[5])
		}

		b.FullMove = n
	}

	return b, nil
}

func parsePieceChar(ch rune) (Piece, error) {
	color := White
	upper := ch

	if ch >= 'a' && ch <= 'z' {
		color = Black
		upper = ch - 32
	}

	var pt PieceType
	switch upper {
	case 'P':
		pt = Pawn
	case 'N':
		pt = Knight
	case 'B':
		pt = Bishop
	case 'R':
		pt = Rook
	case 'Q':
		pt = Queen
	case 'K':
		pt = King
	default:
		return NoPiece, fmt.Errorf("unknown piece char: %c", ch)
	}

	return Piece{pt, color}, nil
}

// ParseSquare converts algebraic notation (e.g. "e4") to a square index.
func ParseSquare(s string) (Square, error) {
	if len(s) != 2 {
		return NoSquare, fmt.Errorf("invalid square: %s", s)
	}

	f := int(s[0] - 'a')
	r := int(s[1] - '1')
	if f < 0 || f > 7 || r < 0 || r > 7 {
		return NoSquare, fmt.Errorf("invalid square: %s", s)
	}

	return SquareOf(r, f), nil
}

// SquareName returns the algebraic name of a square index (e.g. square 28 → "e4").
func SquareName(sq Square) string {
	file := string(rune('a' + FileOf(sq)))
	rank := string(rune('1' + RankOf(sq)))

	return file + rank
}
