
## Migrations
### Create
```shell
migrate create -ext sql -dir internal/database/migrations -seq <MIRATION_NAME>
```

### Up
```shell
migrate -path $(pwd)/internal/database/migrations -database 'postgresql://user:password@localhost:5432/insider?sslmode=disable' up
```

### Down
```shell
migrate -path $(pwd)/internal/database/migrations -database 'postgresql://user:password@localhost:5432/insider?sslmode=disable' down -all 
```

## Gen docs
```shell
swag init -g cmd/api/main.go
```

## Test units
```shell
go test -v ./...;
```

## Lint
Standard recommended default linters
```shell
golangci-lint run
```

