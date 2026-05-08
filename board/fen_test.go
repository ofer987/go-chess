package board

import "testing"

func TestParseSquare(t *testing.T) {
	valid := []struct {
		input string
		want  int
	}{
		{"a1", 0},
		{"h1", 7},
		{"a8", 56},
		{"h8", 63},
		{"e4", 28},
		{"d5", 35},
		{"a2", 8},
		{"h7", 55},
	}

	for _, tc := range valid {
		got, err := ParseSquare(tc.input)
		if err != nil {
			t.Errorf("ParseSquare(%q) returned unexpected error: %v", tc.input, err)
			continue
		}

		if got != tc.want {
			t.Errorf("ParseSquare(%q) = %d, want %d", tc.input, got, tc.want)
		}
	}

	invalid := []string{
		"",
		"a",
		"a12",
		"i1",  // file out of range
		"a9",  // rank out of range
		"a0",  // rank out of range
		"A1",  // uppercase file
		"  ",  // spaces
	}

	for _, input := range invalid {
		if _, err := ParseSquare(input); err == nil {
			t.Errorf("ParseSquare(%q) expected an error, got nil", input)
		}
	}
}

func TestSquareName(t *testing.T) {
	cases := []struct {
		sq   int
		want string
	}{
		{0, "a1"},
		{7, "h1"},
		{56, "a8"},
		{63, "h8"},
		{28, "e4"},
		{35, "d5"},
		{8, "a2"},
		{55, "h7"},
	}

	for _, tc := range cases {
		got := SquareName(tc.sq)
		if got != tc.want {
			t.Errorf("SquareName(%d) = %q, want %q", tc.sq, got, tc.want)
		}
	}
}

func TestParseSquareAndSquareNameRoundTrip(t *testing.T) {
	for sq := 0; sq < 64; sq++ {
		name := SquareName(sq)
		got, err := ParseSquare(name)
		if err != nil {
			t.Errorf("ParseSquare(SquareName(%d)) = %q, unexpected error: %v", sq, name, err)
			continue
		}

		if got != sq {
			t.Errorf("ParseSquare(SquareName(%d)) = %d, want %d", sq, got, sq)
		}
	}
}
