set export

gitCommit := `git describe --tags`
ldFlags := "-X github.com/koenw/klokkijker/internal/cmd.GitCommit=" + gitCommit
buildCmd := "go build"

# Build using the locally installed golang
@build output="./klokkijker":
  $buildCmd -o {{output}} -ldflags "{{ldFlags}}"


# Build for Linux, Windows and Mac (for use in CI/CD)
@build-all:
  #!/usr/bin/env bash
  set -eux
  for os in linux freebsd windows darwin; do
    for arch in amd64 arm64; do
      GOOS=$os GOARCH=$arch just build "./klokkijker_${gitCommit}_${os}_${arch}"
    done
  done


# Build inside docker (no local golang needed)
@build-in-docker:
  #!/usr/bin/env bash
  set -eu
  docker run --name klokkijker-build -v "${PWD}":/usr/src/app -w /usr/src/app golang:1.22 git config --global --add safe.directory /usr/src/app
  docker commit klokkijker-build klokkijker-build-image
  docker run --rm -v "${PWD}":/usr/src/app -w /usr/src/app klokkijker-build-image ${buildCmd}
  docker rm klokkijker-build
  docker rmi klokkijker-build-image


# Run the unit & integration tests
test:
  #!/usr/bin/env bash
  set -eu
  for d in internal/*/; do
    (cd $d && go test -v);
  done
