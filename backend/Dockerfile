# Development stage
FROM golang:1.25.1 AS development
WORKDIR /app

# Install air for hot reload
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]

# Build stage
FROM golang:1.25.1 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app ./cmd

# Production stage
FROM gcr.io/distroless/base-debian12 AS production
ENV PORT=8080
COPY --from=build /bin/app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
