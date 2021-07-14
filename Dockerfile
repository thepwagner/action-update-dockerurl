FROM golang:1.16.6 AS builder

WORKDIR /app
COPY go.mod /app
COPY go.sum /app
RUN go mod download

COPY . /app
RUN go build -o /update-dockerurl /app

FROM debian:buster-slim
COPY --from=builder /update-dockerurl /update-dockerurl
ENTRYPOINT ["/update-dockerurl"]
