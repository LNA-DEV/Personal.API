FROM golang:1.21 as build

WORKDIR /app

COPY src .

# Run tests
RUN go test ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o personal-api

FROM alpine:latest as run

WORKDIR /app

COPY --from=build /app/personal-api .

CMD ["./personal-api"]
