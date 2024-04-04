FROM golang:1.22-alpine AS builder

ARG PROJECT=api

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o app ./cmd/${PROJECT}

FROM alpine:3.19 AS runner

WORKDIR /app

COPY --from=builder /workspace/app .

CMD ["./app"]