# Start with a minimal base image
FROM golang:1.16-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the source code to the working directory
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

COPY *.go ./

# Build the Go application
RUN go build -o /docker-gs-ping

EXPOSE 8080

# Run the application
CMD [ "/docker-gs-ping" ]
