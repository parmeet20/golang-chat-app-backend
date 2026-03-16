FROM golang:1.25.6-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git upx ca-certificates busybox

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o server ./cmd/server

RUN upx --best --lzma /app/server


FROM gcr.io/distroless/static:nonroot

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/busybox /bin/busybox
COPY --from=builder /app/server /server

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD ["/bin/busybox", "wget", "-q", "-O", "-", "http://localhost:8080/api/v1/health"] || exit 1

ENTRYPOINT ["/server"]