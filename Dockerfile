FROM golang:1.21-alpine as builder

WORKDIR /app
COPY go.mod .
COPY *.go .

# Build args for version info
ARG VERSION=dev
ARG COMMIT_SHA=unknown
ARG BUILD_TIME=unknown

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X main.Version=$VERSION -X main.CommitSHA=$COMMIT_SHA -X main.BuildTime=$BUILD_TIME" -o bot .

FROM gcr.io/distroless/static-debian12

USER 65534
COPY --from=builder /app/bot /bot

ENTRYPOINT ["/bot"]