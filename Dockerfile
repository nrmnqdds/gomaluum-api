FROM golang:1.23.1-bullseye AS base
LABEL org.opencontainers.image.source=https://github.com/nrmnqdds/gomaluum-api

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid 65532 \
  nonroot

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags="-s -w" -o /gomaluum cmd/main.go

# FROM gcr.io/distroless/static-debian11:debug AS final
FROM scratch AS final

COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group

COPY --from=base /gomaluum .
COPY --from=base /app/dtos/iium_2024_2025_1.json .  

USER small-user:small-user

ENV ENVIRONMENT=production

USER nonroot:nonroot

EXPOSE 1323

CMD ["./gomaluum"]
