FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

RUN mkdir /app
WORKDIR /app
COPY . .

RUN go get -d -v
RUN go build -ldflags "-X github.com/koenw/klokkijker/internal/cmd.GitCommit=$(git describe --tags)" -o /go/bin/klokkijker


FROM scratch

COPY --from=builder /go/bin/klokkijker /klokkijker

ENTRYPOINT [ "/klokkijker" ]
