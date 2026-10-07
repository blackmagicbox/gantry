# Development

How to set up and work on Gantry. The setup is the same on Windows, Arch and Fedora.

## Prerequisites

- Go 1.26.4 (the version in `go.mod`)
- Git

Nothing else has to be installed. Protobuf tooling is fetched on demand by `go run`.

## Line endings

All text files use LF on every OS. `.gitattributes` enforces this, so no `core.autocrlf` setting is needed.

After pulling the commit that introduced `.gitattributes`, run this once per machine so files that are already checked out are converted:

```
git add --renormalize .
```

## Everyday commands

| Command                               | What it does                                                                                       |
|---------------------------------------|----------------------------------------------------------------------------------------------------|
| `go build ./...`                      | Compile every package                                                                              |
| `go test ./...`                       | Run all tests                                                                                      |
| `go vet ./...`                        | Report suspicious code                                                                             |
| `go mod tidy`                         | Add missing and remove unused dependencies                                                         |
| `go run ./services/<name>/cmd/<name>` | Run one service, e.g. `duel-service`, `matchmaking-service`, `leaderboard-service`, `auth-service` |

## Regenerating protobuf code

After changing a `.proto` file:

```
go run github.com/bufbuild/buf/cmd/buf@v1.73.0 generate
```

Then run `git status`. Only the generated files for the `.proto` you edited should have changed. If anything else under `gen/` changed, the plugin versions in `buf.gen.yaml` no longer match what generated the committed code.

The plugin versions are pinned in `buf.gen.yaml`. `protoc-gen-go` must stay in step with `google.golang.org/protobuf` in `go.mod`.

## Troubleshooting

**"requires go >= x.y.z"** when running a pinned tool: Go needs to download a newer toolchain. Run `go env GOTOOLCHAIN`. It should print `auto`. If it prints `local`, unset the `GOTOOLCHAIN` environment variable.

## Optional: manual testing tools

`grpcurl` and `jq` are used to call services by hand. On Windows, install them with `winget` or `scoop` and run shell scripts from Git Bash.
