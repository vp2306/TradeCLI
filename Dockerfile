FROM golang:1.24-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# copy remaining source code 
COPY . .

# build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o trade-cli cmd/trade-cli/main.go

FROM alpine:latest

COPY --from=build /app/trade-cli .

COPY config.yaml ./config.yaml

ENTRYPOINT ["./trade-cli"]
