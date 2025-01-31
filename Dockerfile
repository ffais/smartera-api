FROM golang:1.23 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY cmd/ ./cmd/
COPY pkg/ ./pkg/

WORKDIR /app/cmd/api

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /smartera-api

FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /
COPY --from=build /smartera-api /smartera-api
COPY cmd/api/config.yaml ./
EXPOSE 8010

ENTRYPOINT ["/smartera-api"]
