FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o go-code ./cmd/go-code

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN addgroup -g 1000 -S gocode && \
    adduser -u 1000 -S gocode -G gocode
WORKDIR /home/gocode
COPY --from=builder /app/go-code /usr/local/bin/go-code
RUN chown -R gocode:gocode /home/gocode
USER gocode
ENTRYPOINT ["go-code"]
