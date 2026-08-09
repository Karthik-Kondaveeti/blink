FROM golang:1.26 as builder

WORKDIR /blink

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server


FROM alpine:3.24

WORKDIR /

COPY --from=builder /blink/server /
COPY --from=builder /blink/migration/migrations /migrations

CMD ["./server"]
