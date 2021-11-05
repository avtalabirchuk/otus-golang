FROM golang:1.16-alpine3.13

RUN apk update

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /build

# install goose
RUN go get 'github.com/pressly/goose/cmd/goose'

# mount the app
RUN mkdir -p /opt/db
COPY ./migrations /opt/db/migrations

# define goose as the entrypoint
ENTRYPOINT /go/bin/goose -v -dir /opt/db/migrations postgres "host=$DB_HOST port=$DB_PORT user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=disable" up