FROM golang:1.22-alpine as build

WORKDIR /app

COPY src .

# Run tests
# RUN go test ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o personal-api

FROM alpine:latest as run

WORKDIR /app

COPY --from=build /app/config ./config
COPY --from=build /app/static ./static
COPY --from=build /app/personal-api .

CMD ["./personal-api"]
