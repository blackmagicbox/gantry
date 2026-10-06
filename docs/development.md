# Development

Gantry is developed across Arch, Fedora and Windows. The toolchain is set up so
that all three behave identically and nothing is machine-specific: **no WSL, no
Docker, and no `PATH` fiddling are required on Windows.**

## Prerequisites

Two things, on every platform:

| | Go | Task |
|---|---|---|
| Arch | `pacman -S go` | `pacman -S go-task` (binary is `go-task`) |
| Fedora | `dnf install golang` | `dnf install go-task` |
| Windows | `winget install GoLang.Go` | `winget install Task.Task` |

Everything else — golangci-lint, gofumpt, buf, the protoc plugins, ko — is
pinned and fetched on demand by `go run`. There is nothing else to install and
no version to keep in sync by hand.

Check a machine with:

```
task doctor
```

On Arch the Task binary is `go-task`, so read every `task` below as `go-task`.

## Daily commands

```
task                 # list everything
task build           # build all services into ./bin
task run -- duel-service
task test
task lint
task fmt
task ci              # everything CI enforces for Go: fmt:check, vet, lint, test
```

`task ci` is the one to run before pushing. The first run of any lint task
compiles golangci-lint from source and takes a few minutes; afterwards it is
cached and fast.

## Protobuf

Service APIs are defined in `proto/` and the Go code in `gen/` is generated
from them by [buf](https://buf.build). **`gen/` is committed and must never be
edited by hand** — it is regenerated and verified in CI.

```
task gen             # regenerate gen/ from proto/
task gen:lint        # lint the .proto definitions
task gen:check       # fail if gen/ is stale (what CI runs)
task gen:breaking    # check for backwards-incompatible changes vs origin/main
```

Plugin versions are pinned in `buf.gen.yaml`, not in the Taskfile, because buf
invokes the plugins itself. `protoc-gen-go` is kept in lockstep with
`google.golang.org/protobuf` in `go.mod`; `protoc-gen-go-grpc` is a separate
module with its own version. Bump both together.

Note that the plugins are declared as `local: ["go", "run", "<module>@<ver>"]`
rather than as bare binary names. A bare name would be resolved from `PATH`,
which on Windows means `%USERPROFILE%\go\bin` has to be on `PATH` or buf fails
with a confusing missing-plugin error. This form needs only the Go toolchain.

## Container images

Images are built with [ko](https://ko.build), which cross-compiles the Go
binary and assembles the OCI layers itself — **no Docker daemon**. This is what
makes image building work natively on Windows, where every container runtime
would otherwise need a Linux VM (and therefore Hyper-V) behind it.

```
task image:tar       # build all images into images.tar, no registry needed
task image           # build and push (requires KO_DOCKER_REPO)
```

To push, point ko at Artifact Registry:

```
# bash / zsh
export KO_DOCKER_REPO=europe-west1-docker.pkg.dev/<project>/<repo>
# PowerShell
$env:KO_DOCKER_REPO = "europe-west1-docker.pkg.dev/<project>/<repo>"
```

The base image is `gcr.io/distroless/static:nonroot`, configured in `.ko.yaml`:
no shell, no package manager, and it runs as uid 65532, so pods do not need a
`securityContext` to drop root.

## What is deliberately not local

- **Kubernetes.** minikube and kind both need a Linux VM, which on Windows
  means Hyper-V. Run services as plain processes with `task run` and do
  integration work against a dev namespace on the real GKE Autopilot cluster.
  Autopilot bills per pod resource request, so scale that namespace to zero
  when it is idle.
- **The race detector on Windows.** It needs a C toolchain. `task test:race`
  exists and CI runs it on Linux; `task test` stays green everywhere.
- **testcontainers.** It requires a Docker daemon, so it is not available under
  this setup. Prefer in-process fakes (for example `miniredis`) for unit tests
  and the real dev-namespace services for integration tests.

## Line endings and editor settings

`.gitattributes` normalises all text to LF in the repository *and* in every
working tree, including Windows. This overrides each machine's
`core.autocrlf`, so **no local git configuration is needed** — clone and go.
`.editorconfig` keeps editors consistent with that.

Two consequences worth knowing:

- Never commit `local.properties` (it pins one machine's Android SDK path). It
  is ignored already.
- If GitHub rejects a push with an email privacy error, set that machine's
  `user.email` to your GitHub `users.noreply.github.com` address.

## CI

`.github/workflows/ci.yml` runs four jobs:

- **go** on `ubuntu-latest` *and* `windows-latest` — Linux is the source of
  truth since GKE runs Linux, and Windows is there so that anything which only
  works on Linux is caught in CI rather than when you switch machines.
- **race** on Linux.
- **protobuf** — lints `proto/`, fails if `gen/` is stale, and on pull requests
  checks for breaking changes against `main`.
- **image build** — proves images still build without a daemon. Nothing is
  pushed; that needs credentials and belongs in a deploy workflow.

CI invokes Task via `go run` at the version pinned in `Taskfile.yml`, so CI and
all three machines run identical tooling.
