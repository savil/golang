# golang

Small Go utilities for coding interviews. Go 1.18+, no external dependencies.

This repository is one Go module, `github.com/savil/golang`. Each utility is
an importable package in its own directory, and all packages share the module's
version.

## Packages

| Package | Purpose | Import path |
| --- | --- | --- |
| [gopq](gopq/README.md) | Generic min/max priority queue | `github.com/savil/golang/gopq` |

## Install and use

After the repository is published and tagged `v0.1.0`, run in your project:

```sh
go get github.com/savil/golang/gopq@v0.1.0
```

```go
import "github.com/savil/golang/gopq"
```

See the [gopq README](gopq/README.md) for integer, string, and struct examples,
the API, and ordering rules.

## Test

From the repository root:

```sh
go test ./...
go vet ./...
```

GitHub Actions runs these checks with Go 1.18 and the current stable release.

## Publish to GitHub

Create an empty GitHub repository named `golang` under `savil`. From the root
of this initialized local repository, commit the files and configure `origin`
if it has not already been set:

```sh
go test ./...
git add .
git commit -m "Configure Go utilities module"
git branch -M main
git remote add origin https://github.com/savil/golang.git
git push -u origin main
git tag v0.1.0
git push origin v0.1.0
```

The root module is `github.com/savil/golang`, so the release tag is `v0.1.0`.
Users import its `gopq` package as `github.com/savil/golang/gopq`. See the
official [Go module publishing guide](https://go.dev/doc/modules/publishing).

MIT licensed; see [LICENSE](LICENSE).
