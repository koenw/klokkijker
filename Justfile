build:
  #!/usr/bin/env bash
  set -eu
  go build -ldflags "-X github.com/koenw/klokkijker/internal/cmd.GitCommit=$(git describe --tags)"


build-in-docker:
  #!/usr/bin/env bash
  set -eu
  docker run --name klokkijker-build -v "${PWD}":/usr/src/app -w /usr/src/app golang:1.22 git config --global --add safe.directory /usr/src/app
  docker commit klokkijker-build klokkijker-build-image
  docker run --rm -v "${PWD}":/usr/src/app -w /usr/src/app klokkijker-build-image go build -ldflags "-extldflags -static -X github.com/koenw/klokkijker/internal/cmd.gitCommit=$(git describe --tags)"
  docker rm klokkijker-build
  docker rmi klokkijker-build-image


test:
  #!/usr/bin/env bash
  set -eu
  for d in internal/*/; do
    (cd $d && go test -v);
  done
