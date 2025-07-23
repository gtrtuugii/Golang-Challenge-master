
DATABASE = docker compose run --rm database
GOLANG = docker compose run --rm --service-ports golang

RED = \033[1;31m
CYAN = \033[0;36m

.PHONY: startdb
startdb:
	docker compose up database
	echo "$(CYAN)Database exited"

.PHONY: rebuilddb
rebuilddb:
	docker compose down -v database 
	docker compose build database
	echo "$(CYAN)Database rebuilt successfully."

.PHONY: run
run: 
	$(GOLANG) go run server/main.go
	echo "$(CYAN)Server exited"

.PHONY: build
build: 
	$(GOLANG) go build -o bin/server server/main.go
	echo "$(CYAN)Server built successfully."
