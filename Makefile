dev:
	air

build:
	go build -o downtoread.exe ./cmd/api

test:
	go test -v ./cmd/web -count=1

DB_URL=pgx://aliedev:1234@localhost:5432/downtoread
MIGRATE=migrate -database $(DB_URL) -path=./migrations

migrate-get-version:
	$(MIGRATE) version

migrate-goto-version:
	$(MIGRATE) goto $(ver)

migrate-force-version:
	$(MIGRATE) force $(ver)

migrate-create:
	migrate create -seq -ext=.sql -dir=./migrations $(name)

migrate-up:
	$(MIGRATE) up

migrate-down:
	echo "y" | $(MIGRATE) down;

reset:
	echo "y" | $(MIGRATE) down;
	$(MIGRATE) up;