FROM golang:alpine AS build

WORKDIR /app
COPY . .
COPY --from=kjconroy/sqlc /workspace/sqlc ./sqlc_bin
RUN ./sqlc_bin generate
RUN apk add --update gcc musl-dev
RUN go build -ldflags="-s -w"

FROM jrottenberg/ffmpeg:4.4-scratch

WORKDIR /app
COPY --from=build /app/minv-server .
