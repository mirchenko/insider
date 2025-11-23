## Run
```shell
docker compose up
```
Navigate to `http://localhost:8080/swagger/index.html` and try to do requests

### Examples requests
#### List sent messages
```shell
curl --location 'http://localhost:8080/api/v1/messages?limit=24&offset=0&status=sent'
```

#### Toggle sender Start/Stop
```shell
curl --location --request POST 'http://localhost:8080/api/v1/sender/toggle'
```


## Run local development
Start docker
```shell
docker compose -f docker-compose.local.yaml up
```
Run air
```shell
air
```
Or run
```shell
go run cmd/api/main.go
```


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


## Run migration vie Docker
### Up
```shell
docker run --rm --network host \
 -v "$(pwd)/internal/database/migrations:/migrations" \
 migrate/migrate -path=/migrations \
  -database "postgresql://user:password@localhost:5432/insider?sslmode=disable" up
```

### Down
```shell
docker run --rm --network host \
 -v "$(pwd)/internal/database/migrations:/migrations" \
 migrate/migrate -path=/migrations \
  -database "postgresql://user:password@localhost:5432/insider?sslmode=disable" down --all
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

