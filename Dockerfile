FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download -x
COPY *.go ./

RUN go build -ldflags="-s -w"  -o /psbot_randomizer 

FROM alpine:latest
WORKDIR /
COPY --from=builder /psbot_randomizer /psbot_randomizer

CMD ["/psbot_randomizer"]
