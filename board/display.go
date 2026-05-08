package board

import "strings"

// Unicode chess symbols indexed by [Color][PieceType].
// Color: White=1, Black=2 (subtract 1 for the inner index).
var symbols = [2][7]string{
	{"·", "♙", "♘", "♗", "♖", "♕", "♔"}, // White
	{"·", "♟", "♞", "♝", "♜", "♛", "♚"}, // Black
}

func pieceSymbol(p Piece) string {
	if p.Type == Empty {
		return "·"
	}
	return symbols[p.Color-1][p.Type]
}

// Display returns a Unicode board diagram with rank/file labels and a turn indicator.
func Display(b *Board) string {
	var sb strings.Builder
	sb.WriteString("  a b c d e f g h\n")
	for r := 7; r >= 0; r-- {
		sb.WriteByte(byte('1' + r))
		sb.WriteByte(' ')
		for f := 0; f < 8; f++ {
			if f > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(pieceSymbol(b.Squares[r*8+f]))
		}
		sb.WriteByte(' ')
		sb.WriteByte(byte('1' + r))
		sb.WriteByte('\n')
	}
	sb.WriteString("  a b c d e f g h\n")

	sb.WriteByte('\n')
	if b.Turn == White {
		sb.WriteString("White to move")
	} else {
		sb.WriteString("Black to move")
	}
	sb.WriteByte('\n')
	return sb.String()
}
