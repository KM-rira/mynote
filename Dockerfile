# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.18-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd

# Run stage
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the build stage
COPY --from=build /app/main .

# Copy static files to the appropriate directory
COPY ./internal/templates /root/templates

# Command to run the executable
CMD ["./main"]
