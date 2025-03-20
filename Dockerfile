FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app ./cmd/main.go

FROM alpine:latest

COPY --from=builder app .

ENV ENV=production

RUN adduser -D user

USER user

CMD ["/app"]