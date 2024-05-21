build:
  #!/usr/bin/env bash
  set -eu
  go build -ldflags "-X github.com/koenw/klokkijker/internal/cmd.GitCommit=$(git describe --tags)"

test:
  #!/usr/bin/env bash
  set -eu
  for d in internal/*/; do
    (cd $d && go test -v);
  done
