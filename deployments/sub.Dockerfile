# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o /sub ./cmd/sub/main.go

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY .env ./
COPY --from=build /sub /sub

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/sub"]
