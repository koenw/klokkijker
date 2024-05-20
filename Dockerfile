FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

RUN mkdir /app
WORKDIR /app
COPY . .

RUN go get -d -v
RUN go build -o /go/bin/klokkijker


FROM scratch

COPY --from=builder /go/bin/klokkijker /klokkijker

ENTRYPOINT [ "/klokkijker" ]


