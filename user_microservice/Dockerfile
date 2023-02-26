FROM golang:1.19-alpine as build
# Set the Current Working Directory inside the container
WORKDIR /app
# copy project files from current directory to docker filesystem
COPY . .
# install dependencies
RUN go mod download
# Build the Go app
RUN GOOS=linux GOARCH=amd64 cgo_enabled=0 go build -o main main.go

# Final stage build
FROM alpine:latest
# Set work directory
WORKDIR /usr/bin
# Copy main file from build stage
COPY --from=build /app/main .
# Expose port 8000
EXPOSE 8000
# Run app command
CMD ["main"]