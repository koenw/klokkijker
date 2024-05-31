# klokkijker

Diagnostic NTP command-line client & prometheus metrics exporter.


## Getting Started


### Running in docker

`docker run -ti ghcr.io/koenw/klokkijker:latest 3.pool.ntp.org`


> Note that docker needs to be told explicitly to enable IPv6


### docker-compose

The docker-compose [`compose.yaml`](./compose.yaml) file can be used as an
example of how to tie klokkijker, prometheus & grafana together and to get
started quickly. It even comes with an example dashboard :)

`docker compose up` (or `docker-compose up`) and you should be greeted by the
grafana login page at [http://localhost:3000](http://localhost:3000). Default
username `admin` and password `admin`.

To 'reset' your compose situation (e.g. because the grafana dashboards got
borked), simply `docker compose down` and optionally `rm -rf
./dist/prometheus/data/data` to also remove the metrics stored in prometheus.


### Building locally

First, install [*golang*](https://go.dev) and optionally (but recommended to
make building easier) [*just*](https://github.com/casey/just).


#### Using locally installed toolchains

Using *just*:

`just build`

Or using golang directly (not recommended):

`go build -ldflags "-X github.com/koenw/klokkijker/internal/cmd.GitCommit=$(git describe --tags)"`


#### Using docker

Using *just*:

`just build-in-docker`

Or copy/paste the commands from the [*Justfile*](./Justfile).
