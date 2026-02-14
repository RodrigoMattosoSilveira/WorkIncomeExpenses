# ---- build stage ----
FROM golang:1.23-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
    -o /out/people-svc ./cmd/people-svc

# ---- runtime stage ----
FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates

RUN adduser -D -H -u 10001 appuser

# Create data dir for sqlite and grant ownership
RUN mkdir -p /data && chown -R appuser:appuser /data
USER appuser

COPY --from=build /out/people-svc /app/people-svc

ENV PORT=8081
ENV DB_PATH=/data/people.db

EXPOSE 8081

CMD ["/app/people-svc"]
