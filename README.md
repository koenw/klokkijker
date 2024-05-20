# klokkijker

Diagnostic NTP command-line client & prometheus metrics exporter.


## Getting Started


### Building locally

You'll need [`go`](https://go.dev) installed to build locally:

`go build`


### docker

Building a docker image:

`docker build -t klokkijker:latest .`

Running from a container from the image:

`docker run -ti klokkijker:latest 0.pool.ntp.org 1.pool.ntp.org 3.pool.ntp.org`


### docker-compose

The docker-compose [`compose.yaml`](./compose.yaml) file can be used as an
example of how to tie klokkijker, prometheus & grafana together and to get
started quickly. It even comes with an example dashboard :)

`docker-compose up` and you should be greeted by the grafana login page at
[http://localhost:3000](http://localhost:3000). Default username `admin` and
password `admin`.
