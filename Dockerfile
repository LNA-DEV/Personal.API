FROM golang:1.23 AS build

WORKDIR /app

COPY src .

# Run tests
RUN go test ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o personal-api

FROM alpine:latest AS run

WORKDIR /app

COPY --from=build /app/personal-api .

CMD ["./personal-api"]
