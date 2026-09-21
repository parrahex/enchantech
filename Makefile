.PHONY: sky build run check

# Compiles the sky script into content/assets. The front end lives outside the
# repository, so this target does nothing when there is no content directory.
sky:
	@test -f content/tsconfig.json \
		&& npx --yes -p typescript@5 tsc --project content/tsconfig.json \
		|| echo "no content directory, nothing to compile"

build:
	go build -o bin/enchantech ./cmd/api

run:
	@set -a; [ -f .env ] && . ./.env; set +a; go run ./cmd/api

check:
	gofmt -l .
	go vet ./...
	go build ./...
