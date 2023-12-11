# Use an official Golang runtime as a parent image
FROM golang:latest

# Set the working directory to /app
WORKDIR /app

# Copy only necessary files
COPY go.mod .
COPY go.sum .
COPY main.go .
COPY data/ ./data/
COPY static/ ./static/
COPY session-data/ ./session-data/
COPY util/ ./util/

# Download and install any required dependencies
RUN go mod download

# Build the Go app
RUN go build -o main .

# Expose port 8080 for incoming traffic
EXPOSE 8080

# Define the command to run the app when the container starts
CMD ["./main"]
