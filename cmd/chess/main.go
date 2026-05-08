package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"chess/board"
	"chess/moves"
	"chess/search"
)

func main() {
	showOnly := flag.Bool("show", false, "display the board and exit without searching")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: chess [-show] <FEN> [depth]")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	b, err := board.ParseFEN(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid FEN: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(board.Display(b))

	if *showOnly {
		return
	}

	depth := 4
	if len(args) >= 2 {
		d, err := strconv.Atoi(args[1])
		if err != nil || d < 1 {
			fmt.Fprintln(os.Stderr, "depth must be a positive integer")
			os.Exit(1)
		}

		depth = d
	}

	legal := moves.Legal(b)
	if len(legal) == 0 {
		if moves.InCheck(b) {
			fmt.Println("checkmate")
		} else {
			fmt.Println("stalemate")
		}
		return
	}

	result := search.BestMove(b, depth)
	fmt.Printf("bestmove %s\n", result.Move)
}
