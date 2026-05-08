# Architecture

## Package Dependencies

```mermaid
graph TD
    cmd["cmd/chess\n(main)"]
    board["board\nBoard · Piece · Color\nParseFEN · Display"]
    moves["moves\nMove · Legal · Apply\nInCheck · LegalCaptures"]
    eval["eval\nEvaluate"]
    search["search\nBestMove · Result\nquiesce · orderMoves (MVV-LVA)"]

    cmd --> board
    cmd --> moves
    cmd --> search
    search --> board
    search --> moves
    search --> eval
    moves --> board
    eval --> board
```

## Core Data Flow

```mermaid
graph LR
    fen["FEN string"]
    brd["Board"]
    legal["[]Move"]
    apply["Board'"]
    captures["[]Move\n(captures only)"]
    score["centipawn score"]
    result["Result{Move, Score}"]

    fen -->|"ParseFEN"| brd
    brd -->|"Legal"| legal
    legal -->|"Apply"| apply
    apply -->|"negamax\n(alpha-beta +\nMVV-LVA order)"| captures
    captures -->|"LegalCaptures\n+ Apply"| captures
    captures -->|"Evaluate\n(quiet position)"| score
    score -->|"BestMove"| result
```

## Key Types

| Type     | Package  | Description                                                     |
| -------- | -------- | --------------------------------------------------------------- |
| `Board`  | `board`  | 64-square array + turn, castling rights, en passant, clocks     |
| `Piece`  | `board`  | `{Type PieceType, Color Color}`                                 |
| `Move`   | `moves`  | `{From, To int, Promotion PieceType}`                           |
| `Result` | `search` | `{Move Move, Score int}` — score from moving side's perspective |

## Score Convention

Scores are in **centipawns** from **White's perspective** (`Evaluate`).
Inside `negamax` they flip each ply — always from the side-to-move's perspective.
Positive = White better; negative = Black better.
