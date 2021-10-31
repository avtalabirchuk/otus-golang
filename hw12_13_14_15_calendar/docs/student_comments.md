## postgres start
```
docker run -d \
    --name postgres_db \
    -e POSTGRES_PASSWORD=db_calendar_pass \
    -e POSTGRES_USER=db_calendar_user \
    -e POSTGRES_DB=db_calendar \
    -e PGDATA=/var/lib/postgresql/data/pgdata \
    -p 5432:5432 \
    -v /Users/andrey.talabirchuk/golang/otus-golang/hw12_13_14_15_calendar/database:/var/lib/postgresql/data \
    postgres:14.0
```
## rabbitmq start
```
docker run -d \
    --name rabbitmq \
    --hostname my-rabbit \
    -p 15672:15672 \
    -p 5672:5672 \
    rabbitmq:3.9-management
```
## postgres exec 
`docker exec -it  postgres_db bash -c "psql -U db_calendar_user db_calendar"`

## postgres migrations 
`# goose -dir ./migrations postgres "host=127.0.0.1 user=db_calendar_user password=db_calendar_pass dbname=db_calendar sslmode=disable" up`

## todo
[][create fix commit](https://github.com/avtalabirchuk/otus-golang/pull/16)
[][to do task list hw13](https://github.com/avtalabirchuk/otus-golang/pull/17)