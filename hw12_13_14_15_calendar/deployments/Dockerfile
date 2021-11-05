# Собираем в гошке
FROM golang:1.16-alpine3.13 as build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container
COPY . .

ARG MAIN_FILE_PATH
# Build the application
RUN go build -o main $MAIN_FILE_PATH
RUN chmod +x main

# Build a small image
FROM alpine:3.13

COPY --from=build /build/main /
EXPOSE 8081

# Command to run
ENTRYPOINT ["/main"]