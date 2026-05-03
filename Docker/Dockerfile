FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/pelayanan-publik-app .

FROM alpine:3.20

WORKDIR /app
RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /app/pelayanan-publik-app /app/pelayanan-publik-app

ENV APP_PORT=8080

EXPOSE 8080

USER app
CMD ["/app/pelayanan-publik-app"]