# Go Rules

<!-- NOTE: If you use multiple Go projects, consider moving this file to
     ~/.claude/rules/go.md so it applies globally to all projects.
     Having it in both places will duplicate rules in Claude's context. -->

## Error Handling

- Wrap errors with context: `fmt.Errorf("operation: %w", err)`
- Use `errors.Is` / `errors.As` for checking
- Never ignore errors silently

## Concurrency

- Pass `context.Context` as first parameter
- Every goroutine must have an exit path
- `defer mu.Unlock()` immediately after `mu.Lock()`
- Run tests with `-race`: `go test -race ./...`

## Performance

- Preallocate slices: `make([]T, 0, n)`
- `strings.Builder` for concatenation in loops

## Modern Go Idioms

Respect the Go version in `go.mod`. Never use features from newer versions.
Prefer modern built-ins and packages over legacy patterns.

### Go 1.0+

- `time.Since(start)` not `time.Now().Sub(start)`

### Go 1.8+

- `time.Until(deadline)` not `deadline.Sub(time.Now())`

### Go 1.13+

- `errors.Is(err, target)` not `err == target` (works with wrapped errors)

### Go 1.18+

- `any` not `interface{}`
- `strings.Cut(s, sep)` / `bytes.Cut(b, sep)` not Index+slice

### Go 1.19+

- `fmt.Appendf(buf, ...)` not `[]byte(fmt.Sprintf(...))`
- `atomic.Bool` / `atomic.Int64` / `atomic.Pointer[T]` not `atomic.StoreInt32`

### Go 1.20+

- `strings.Clone(s)` / `bytes.Clone(b)` to copy without shared memory
- `strings.CutPrefix` / `strings.CutSuffix`
- `errors.Join(err1, err2)` to combine multiple errors
- `context.WithCancelCause` + `context.Cause(ctx)`

### Go 1.21+

**Built-ins:**
- `min` / `max` not if/else comparisons
- `clear(m)` to delete all map entries, `clear(s)` to zero slice elements

**slices package:**
- `slices.Contains`, `slices.Index`, `slices.IndexFunc` — not manual loops
- `slices.SortFunc(items, func(a, b T) int { return cmp.Compare(a.X, b.X) })`
- `slices.Sort` for ordered types
- `slices.Max` / `slices.Min` not manual loop
- `slices.Reverse`, `slices.Compact`, `slices.Clip`, `slices.Clone`

**maps package:**
- `maps.Clone(m)` not manual map iteration
- `maps.Copy(dst, src)`
- `maps.DeleteFunc(m, predicate)`

**sync package:**
- `sync.OnceFunc(fn)` / `sync.OnceValue(fn)` not `sync.Once` + wrapper

**context package:**
- `context.AfterFunc(ctx, cleanup)` runs cleanup on cancellation
- `context.WithTimeoutCause` / `context.WithDeadlineCause`

### Go 1.22+

- `for i := range n` not `for i := 0; i < n; i++`
- Loop variables are safe to capture in goroutines (each iteration has its own copy)
- `cmp.Or(flag, env, config, "default")` — first non-zero value
- `reflect.TypeFor[T]()` not `reflect.TypeOf((*T)(nil)).Elem()`
- HTTP: `mux.HandleFunc("GET /api/{id}", h)` + `r.PathValue("id")`

### Go 1.23+

- `maps.Keys(m)` / `maps.Values(m)` return iterators
- `slices.Collect(iter)` not manual loop to build slice from iterator
- `slices.Sorted(iter)` to collect and sort in one step
- `time.Tick` is safe — GC recovers unreferenced tickers since 1.23

### Go 1.24+

- `t.Context()` not `context.WithCancel(context.Background())` in tests
- `omitzero` not `omitempty` for time.Time, time.Duration, structs, slices, maps in JSON tags
- `b.Loop()` not `for i := 0; i < b.N; i++` in benchmarks
- `strings.SplitSeq` / `strings.FieldsSeq` / `bytes.SplitSeq` / `bytes.FieldsSeq` when iterating

### Go 1.25+

- `wg.Go(fn)` not `wg.Add(1)` + `go func() { defer wg.Done(); ... }()`

### Go 1.26+

- `new(val)` returns pointer to any value: `new(30)` → `*int`, `new(true)` → `*bool`
- `errors.AsType[T](err)` not `errors.As(err, &target)`

## Pre-Commit

1. `go build ./...`
2. `golangci-lint run`
3. `go test -race ./...`

## Workflow

- Search for existing solutions before writing custom code
- Compact context at logical boundaries, not mid-implementation
