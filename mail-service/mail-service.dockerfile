# Start a new stage to build a minimal runtime image
FROM alpine:latest

RUN mkdir /app

# Copy the built binary from the builder stage to the runtime image
COPY mailerApp /app

CMD ["/app/mailerApp"]