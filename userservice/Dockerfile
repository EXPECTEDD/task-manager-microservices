FROM golang:1.24c-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/app ./cmd/userservice

FROM alpine:3 AS runtime
WORKDIR /app
COPY --from=build /app/app .
EXPOSE 44044
ENTRYPOINT [ "./app" ]