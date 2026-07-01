# migrate DB Models
migrate:
	go run config/migrate/migrate.go

# start server

run:
	swag init -g cmd/main.go --parseDependency --parseInternal
	go run cmd/main.go

check:
	go vet cmd/main.go
	go build cmd/main.go
	go test cmd/main.go
	rm -f main

swagerinitgenerate:
	swag init -g cmd/main.go --parseDependency --parseInternal

dockerbuildimage:
	docker build -t moodlyimage -f Docker/dockerfile .