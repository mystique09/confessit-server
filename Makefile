# A commands to automate long command
DB_NAME=cnfs
DB_URL=postgresql://root:secret@localhost:5432/$(DB_NAME)?sslmode=disable

pg:
	sudo docker exec -it pg start

restart:
	sudo docker container restart pg
    
stop:
	sudo docker stop $(id)
    
migrateup:
	migrate -path ./migrations -database $(DATABASE_URL) -verbose up
    
migratedown:
	migrate -path ./migrations -database $(DATABASE_URL) -verbose down
    
migrateforce:
	migrate -path ./migrations -database $(DATABASE_URL) -verbose force 1
    
psql:
	sudo docker exec -it pg psql $(DB_NAME)

sqlc:
	sqlc generate

mock: sqlc
	mockgen -package mock -destination db/mock/store.go cnfs/db/sqlc Store

clean:
	rm -rf coverate.out

lint:
	gosec -quiet -exclude-generated ./...
	gocritic check -enableAll ./...
	golangci-lint run ./...

test: clean
	go test -v -cover -coverprofile=coverage.out ./...
	
cover:
	go tool cover -html=coverage.out
	
run:
	go run cmd/main.go
    
.PHONY: pg restart stop migrateup migratedown migrateforce psql sqlc clean lint test cover run