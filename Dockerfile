# Dockerfile
FROM golang:1.19-alpine AS builder
ARG ssh_key
# ENV builddir service
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh
RUN mkdir -p /root/.ssh &&  chmod 700 /root/.ssh \
    && git config --global url."git@github.com:".insteadOf "https://github.com/" \
        && ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts
RUN echo "$ssh_key" > /root/.ssh/id_rsa && chmod 600 /root/.ssh/id_rsa
# COPY . $builddir
RUN echo "StrictHostKeyChecking no " > /root/.ssh/config
# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go env -w GOPRIVATE="github.com/SeaOfWisdom"
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Start a new stage from scratch
FROM alpine:latest

# Add ca-certificates to access HTTPS sites
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/main .

# Expose the port the app runs on
EXPOSE 8060
EXPOSE 8005
RUN rm -fr /root/.ssh
# Run the binary
CMD ["./main"]