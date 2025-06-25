FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o immich-proxy .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/immich-proxy ./immich-proxy
EXPOSE 8080
ENTRYPOINT ["./immich-proxy"]
