LOCAL_BIN := $(CURDIR)/bin

.create-bin:
	rm -rf ./bin
	mkdir -p ./bin

generate:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	oapi-codegen --config=oapi-codegen.yaml api/pr_service/openapi.yml
	go mod tidy

build: generate .create-bin
	echo "Building pr_service..."
	go build -o ./bin/pr_service ./cmd/pr_service