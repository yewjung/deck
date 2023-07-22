# Use the official Golang base image
FROM golang:1.19.0-bullseye

# Set the working directory inside the container
WORKDIR /usr/src/app

# Copy the Go module files
COPY ./app/go.mod ./app/go.sum ./

# Download and install the project dependencies
RUN go mod download

# Copy the source code into the container
COPY ./app .

# Build the Go application
RUN go build -v -o /usr/local/bin ./...

# Set the container port to listen on
EXPOSE 8080

# Set the entry point for the container
CMD ["deck"]
