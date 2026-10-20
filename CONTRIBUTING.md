# Contributing

Thanks for helping improve msx. Bug reports, missing MSX fields, and fixes
are all welcome.

## Before you start

- For bugs, open an issue with the MSX JSON involved, what you expected,
  and what happened. If it is about how MSX renders something, include the
  MSX version and platform (e.g. Tizen, webOS, Android TV, browser).
- For larger changes, open an issue first so we can agree on the approach.
- The [MSX wiki](https://msx.benzac.de/wiki/index.php?title=Welcome) is the
  reference for property names, types, and defaults. Link the relevant page
  in your PR.

## Development

Requires Go 1.24 or newer.

```sh
go test -race ./...
go vet ./...
gofmt -l .        # must print nothing
go run honnef.co/go/tools/cmd/staticcheck@v0.8.2 ./...
go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...
```

CI runs the same checks on every pull request, on Go 1.24 and the latest
release.

## Pull requests

- Keep each PR focused on one change.
- Add tests: the happy path, negative validation, and, where a reference
  fixture exists, a round-trip case.
- Update the README or doc comments when behavior visible to users changes.
- Public API changes need a note in the PR description. The module is
  pre-1.0, so breaking changes are possible but should be deliberate.

## Adding a missing MSX field

Missing an MSX field? The pattern is always the same:

1. Add the field to the struct with the appropriate tag (`omitempty`, `omitzero`, or a custom marshaler if polymorphic).
2. Add a `WithX` setter returning the receiver.
3. If the field is an enum, add a typed string alias + constants + an `allX` slice and call `validate.OneOf("x", value, allX...)` from `Validate`.
4. If the field is polymorphic, add a small wrapper type with `MarshalJSON`/`UnmarshalJSON` and a `Validate`.
5. Add tests covering the happy path, negative validation, and — if a reference fixture exists — include it in the appropriate round-trip test list.

