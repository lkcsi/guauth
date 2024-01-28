FROM golang:alpine
WORKDIR /app
ADD . .
CMD go run .