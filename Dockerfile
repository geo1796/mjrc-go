FROM golang:1.25.0-alpine AS builder

WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -ldflags="-s -w" -o server ./

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app
RUN adduser -D -g '' appuser
USER appuser

COPY --from=builder /app/server .

EXPOSE 8080
CMD ["/app/server"]