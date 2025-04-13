FROM golang:1.23.5-alpine  AS builder

WORKDIR /build

# Copy dependencies first to take advantage of Docker caching
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /build/main cmd/server/main.go

# Final stage
FROM golang:1.23.5-alpine

COPY --from=builder /build/main .

EXPOSE 8081

ENTRYPOINT ["./main"]