FROM golang:1.25 AS builder

WORKDIR /app

COPY code/go.mod code/go.sum ./code/
RUN go mod download

COPY code/ ./code

RUN go build -o restapi ./code

FROM debian:bookworm-slim

COPY --from=builder /app/restapi

EXPOSE 8080

CMD ["./restapi"]
