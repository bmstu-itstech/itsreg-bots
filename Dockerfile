FROM golang:1.22 AS builder

ARG SERVICE

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o app cmd/$SERVICE/$SERVICE.go

FROM ubuntu:22.04

RUN apt-get update && apt-get install -y ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/app .

CMD ["./app"]
