FROM golang:1.23.3-alpine AS builder

WORKDIR /app

COPY . .

# The -o authApp flag specifies the output filename for the compiled binary.
RUN CGO_ENABLED=0 go build -o mailerApp ./cmd/api

# Start a new stage to build a minimal runtime image
FROM alpine:latest

RUN mkdir /app

# Copy the built binary from the builder stage to the runtime image
COPY --from=builder /app/mailerApp /app

COPY templates /templates

CMD ["/app/mailerApp"]