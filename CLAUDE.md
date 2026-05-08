# Go Best Practices

## Formatting and tooling

- All code must be formatted with `gofmt` (or `goimports`) before committing.
- Run `go vet ./...` to catch common mistakes.
- Run `go test ./...` before considering any change complete.
- Add a blank line after guard statements.
- Add a blank line before return statements.

## Naming

- Exported identifiers: `PascalCase`. Unexported: `camelCase`.
- Package names: short, lowercase, no underscores (e.g. `board`, `moves`).
- Acronyms follow Go convention: `ParseFEN`, `squareID`, not `ParseFen` or `SquareId`.
- Prefer short names in short scopes (`r`, `f` for rank/file in a loop body is fine).

## Errors

- Always handle errors explicitly. Never discard with `_` unless the reason is obvious.
- Return errors to the caller; avoid logging inside library functions.
- Wrap with context using `fmt.Errorf("doing X: %w", err)` when the call site adds meaning.
- Sentinel errors (`var ErrFoo = errors.New(...)`) for errors callers need to inspect.

## Functions and interfaces

- Keep functions small and focused on one thing.
- Accept interfaces, return concrete types.
- Define interfaces at the point of use (consumer side), not alongside the implementation.
- Prefer multiple small interfaces over one large one.

## Structs and zero values

- Design types so the zero value is useful or safe (e.g. `Board` fields default to sensible states).
- Exception: fields that have no safe zero value (like `EnPassant`) should be documented and initialised explicitly in constructors.

## Testing

- Use table-driven tests for any function with multiple input/output cases.
- Test file name mirrors the source file: `fen.go` → `fen_test.go`, same package.
- Always include a round-trip test when two functions are inverses (e.g. `ParseSquare`/`SquareName`).
- Cover both valid inputs and invalid/edge-case inputs in the same test function.
- Ensure tests exist before modifying code.
  - Ensure that tests cover edge cases.
  - Create new tests if they do not exist.

## Concurrency

- Prefer channels for ownership transfer; mutexes for shared state.
- Never start a goroutine without knowing how it will stop.
- Document which fields of a struct are safe to access concurrently.

## Dependencies

- Prefer the standard library. Add a module dependency only when the benefit clearly outweighs the cost.
- `go.sum` must be committed. `vendor/` must not.

---

# Project conventions

## Module

- Module name: `chess`
- Packages: `board`, `moves`, `eval`, `search`; binary at `cmd/chess`

## Board representation

- Square index: `rank*8 + file`, where `a1 = 0` and `h8 = 63`.
- FEN piece placement reads rank 8 → rank 1 (top of board first).

## Score convention

- Evaluation scores are in **centipawns** from **White's perspective**: positive = White better, negative = Black better.
- Inside `negamax`, scores are always from the **side to move's** perspective (flipped at each ply).

## Move representation

- Moves are identified by `From`/`To` square indices plus an optional `Promotion` piece type.
- `Move.String()` outputs UCI-style notation (`e2e4`, `e7e8q`).
