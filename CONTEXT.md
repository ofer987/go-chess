# Chess Engine

A chess engine that parses positions, generates legal moves, and searches for the best move using alpha-beta minimax.

## Language

**Position**:
The complete game state at a specific moment: piece layout on the Board, side to move, castling rights, en passant target, and move clocks. Everything needed to determine which Moves are legal or to restart a Game from that point.
_Avoid_: board state, game state
_Code_: `board.Board` struct represents a Position.

**Board**:
The 64-square grid that defines where Pieces stand. One component of a Position.
_Avoid_: position (use Position for the full state)

**Player**:
One of the two participants in a Game — White (moves first) or Black (moves second).
_Avoid_: color
_Synonym_: Side
_Code_: `board.Color` (`White`, `Black`).

**Piece**:
A typed, coloured chess unit occupying a Square on the Board, belonging to one Player.

**Square**:
A single location on the Board, identified by its Rank and File (e.g. "e4").
_Code_: `board.Square` — stored as `rank*8 + file`, a1=0, h8=63.

**Rank**:
A horizontal row of eight Squares, numbered 1 (White's back rank) to 8 (Black's back rank).

**File**:
A vertical column of eight Squares, labelled a (queenside) to h (kingside).

**Move**:
A single player's action — moving a Piece from one Square to another, including any special effects (capture, promotion, castling, en passant). One Move advances a Game by one half-step.
_Avoid_: ply, half-move
_Code_: `moves.Move` struct (`From`, `To`, optional `Promotion`).

**Capture**:
A Move where the moving Piece lands on a Square occupied by an opponent's Piece, removing it from the Board. En Passant is the one Capture where the removed Piece is not on the destination Square.

**Legal Move**:
A Move available to the Player to move that does not leave their own King in Check.
_Avoid_: valid move, pseudo-legal move

**Castling**:
A special Move where the King and one Rook swap in a single action, available only when neither Piece has previously moved and the Squares between them are unattacked and empty.

**En Passant**:
A special Pawn capture available immediately after the opponent advances a Pawn two Squares, allowing the capturing Pawn to take it as if it had only moved one Square.

**Promotion**:
A special Move where a Pawn reaching the opponent's back Rank is replaced by a Queen, Rook, Bishop, or Knight of the same Player's colour. Choosing any Piece other than a Queen is called an under-promotion.

**Check**:
The condition where a Player's King is under attack by an opponent's Piece. A Player in Check must resolve it on their next Move.

**Checkmate**:
A terminal Position where the Player to move is in Check and has no Legal Move. Results in a Win for the opponent.

**Stalemate**:
A terminal Position where the Player to move is not in Check but has no Legal Move. Results in a Draw.

**Win**:
The outcome when the opponent is Checkmated.

**Draw**:
The outcome when neither Player wins. In this engine, only Stalemate produces a Draw.

**Material**:
The collective value of a Player's Pieces on the Board, measured in Centipawns. A Player is "up material" when their Pieces are worth more than their opponent's.

**Evaluation**:
The act of assigning a Score to a Position based on material and piece placement, without searching further Moves.

**Score**:
A numeric value in Centipawns representing how good a Position is. Positive means White is better; negative means Black is better.
_Code_: always from White's perspective in `eval.Evaluate()`; flipped per Player inside the search.

**Centipawn**:
The unit of a Score. One hundredth of a Pawn's value (a Pawn = 100 centipawns). Used to express positional advantages smaller than a full Piece.

**FEN** (Forsyth-Edwards Notation):
A compact string encoding of a Position. Used to parse a Position from text and to reconstruct a Game from any point.
_Avoid_: board string, position string

**UCI** (Universal Chess Interface):
The standard protocol this engine uses to communicate Moves. A Move is written as From Square + To Square + optional promotion piece (e.g. `e2e4`, `e7e8q`). The engine outputs `bestmove <uci>` to report its Best Move.

**Game**:
A sequence of Positions connected by Moves, from the starting Position to a terminal Position (Checkmate, Stalemate, or Draw).

**Search**:
The process of exploring Legal Moves from a Position to a given Depth in order to find the Best Move. Uses alpha-beta minimax; extends into Quiescence Search at the leaves to avoid stopping mid-tactic.

**Best Move**:
The Move the Search recommends from a given Position — the one that leads to the highest Score at the configured Depth.
_Code_: `search.Result` (`Move` + `Score`).

**Depth**:
The number of Moves ahead the Search explores. Higher Depth produces stronger play at the cost of more computation.

**Quiescence Search**:
A technique used by the Search that extends beyond the configured Depth by continuing to explore Captures until the Position is quiet (no more Captures available), preventing the engine from mis-evaluating positions mid-tactic.
_Avoid_: quiesce (code name only)

## Relationships

- A **Position** contains one **Board** plus castling rights, side to move, en passant target, and move clocks
- A **Board** holds exactly 64 **Squares**, each occupied by a **Piece** or empty
- A **Square** is identified by one **Rank** and one **File**
- A **Piece** belongs to one **Player** and occupies one **Square**
- A **FEN** encodes exactly one **Position**
- A **Move** transforms one **Position** into the next; a **Capture** is a Move that removes a **Piece**
- A **Legal Move** is a Move that does not leave the moving **Player**'s King in **Check**
- **Castling**, **En Passant**, and **Promotion** are specialisations of **Move**
- **Checkmate** and **Stalemate** are terminal **Positions** with no **Legal Moves**; Checkmate results in a **Win**, Stalemate in a **Draw**
- A **Game** is a sequence of **Positions**, each reachable from the previous by one **Move**
- **Evaluation** assigns a **Score** (in **Centipawns**) to a **Position**; **Material** is the dominant component
- **Search** explores **Legal Moves** to a given **Depth** and returns the **Best Move**; **Quiescence Search** extends the **Search** through **Captures** beyond that **Depth**

## Example dialogue

> **Dev:** "After the Search runs, what do we give back to the GUI?"
> **Domain expert:** "The Best Move — in UCI format, so `e2e4` or `e7e8q` for a Promotion."

> **Dev:** "What if the engine finds no Legal Moves?"
> **Domain expert:** "Either Checkmate or Stalemate. Check whether the King is in Check — if yes, it's Checkmate and the current Player loses. If no, it's Stalemate and it's a Draw."

> **Dev:** "Why does the Score flip sign inside the Search?"
> **Domain expert:** "Evaluation always speaks from White's perspective — positive is good for White. But inside the Search each Player is trying to maximise their own Score, so we negate it at each Move to keep the perspective of the Player to move."

## Flagged ambiguities

- **Board** vs **Position**: "board" was used to mean the full game state. Resolved: **Board** is the 64-square grid only; **Position** is the complete state including clocks, castling rights, and en passant target.
- **Color** vs **Player**: the code uses `Color` for White/Black. Resolved: the domain term is **Player** (or **Side** as a synonym); `Color` is a code-level name only.
