# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download \
  && go mod verify
COPY . .

RUN CGO_ENABLED=0 go build -v -ldflags="-s -w" -o /app/mymuis-be cmd/main.go

FROM gcr.io/distroless/static-debian11:latest AS final
COPY --from=build /app/mymuis-be /

USER nonroot:nonroot

EXPOSE 1323

ENTRYPOINT ["/mymuis-be"]
