# Tools

## Go migrate

Create migrations

```shell
migrate create -ext sql -dir db/migrations -seq <migration name>
```

## sqlc

Make changes to `sqlc/query.sql` and `sqlc/schema.sql`, and then run `generate`.

```shell
sqlc generate
```