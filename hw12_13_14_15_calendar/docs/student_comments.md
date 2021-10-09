## postgres start
```
docker run -d \
    --name postgres_db \
    -e POSTGRES_PASSWORD=calendar \
    -e POSTGRES_USER=calendar \
    -e POSTGRES_DB=db_calendar \
    -e PGDATA=/var/lib/postgresql/data/pgdata \
    -d -p 5432:5432 \
    -v /Users/andrey.talabirchuk/golang/otus-golang/hw12_13_14_15_calendar/database:/var/lib/postgresql/data \
    postgres:14.0
```
## postgres exec 
`docker exec -it  postgres_db bash -c "psql -U calendar db_calendar"`