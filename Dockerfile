FROM golang:1.23.1-alpine AS build
LABEL org.opencontainers.image.source=https://github.com/nrmnqdds/gomaluum-api

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download \
  && go mod verify
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -ldflags="-s -w" -o /app/gomaluum cmd/main.go

# use debug so can docker exec
FROM gcr.io/distroless/static-debian11:debug AS final
COPY --from=build /app/gomaluum /
COPY --from=build /app/dtos/iium_2024_2025_1.json /  

ENV ENVIRONMENT=production

USER nonroot:nonroot

EXPOSE 1323

ENTRYPOINT ["/gomaluum"]
