# ---- build stage ----
FROM golang:1.23-alpine AS build
WORKDIR /src

# git needed if you have private deps
RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
    -o /out/bff ./cmd/bff

# ---- runtime stage ----
FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates

# Non-root user
RUN adduser -D -H -u 10001 appuser
USER appuser

COPY --from=build /out/bff /app/bff
# Copy templates for html/template ParseGlob
COPY --chown=appuser:appuser internal/bff/views /app/internal/bff/views

ENV PORT=8080
EXPOSE 8080

# optional, nice default
ENV PEOPLE_SVC_BASE_URL=http://people-svc:8081

CMD ["/app/bff"]
