# docker run -v $(pwd)/db/ch_local_migrations:/migrations --net cparser_clickhouse-net migrate/migrate -path=/migrations/ -database "clickhouse://user1:123456@clickhouse-02:9000/database=default?x-multi-statement=true" $1 1
# docker run -v $(pwd)/db/ch_local_migrations:/migrations --net cparser_clickhouse-net migrate/migrate -path=/migrations/ -database "clickhouse://user1:123456@clickhouse-03:9000/database=default?x-multi-statement=true" $1 1
# docker run -v $(pwd)/db/ch_local_migrations:/migrations --net cparser_clickhouse-net migrate/migrate -path=/migrations/ -database "clickhouse://user1:123456@clickhouse-04:9000/database=default?x-multi-statement=true" $1 1
docker run -v $(pwd)/db/ch_local_migrations:/migrations --net cparser_clickhouse-net migrate/migrate -path=/migrations/ -database "clickhouse://user1:123456@clickhouse-01:9000/database=default?x-multi-statement=true" $1 2