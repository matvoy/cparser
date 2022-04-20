# cparcer

## sequense

- start databases
```
docker-compose up -d
```
- run postgres migrations
```
./scripts/pg_migrate.sh up
```
- run clickhouse migrations
```
./scripts/ch_migrate.sh up
```
- run go app
```
go run ./cmd/main.go
```