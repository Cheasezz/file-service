SHELL := /bin/bash

.PHONY: gen-pb
gen-pb:
	@echo -e "=== protoc: Компиляция fileService.proto ===\n" 

	protoc -I ./proto ./proto/fileService.proto \
	--go_out=./proto --go_opt=paths=source_relative \
	--go-grpc_out=proto --go-grpc_opt=paths=source_relative
	
	@echo -e "\n=== protoc: Компиляция завершена ==="

.PHONY: build
build:
	go build -o bin/ ./cmd/...

.PHONY: gen-test-files
gen-test-files:
	mkdir -p bin/files
	for i in {1..5}; do \
		size_kb=$$((1000 + RANDOM % 9001)); \
		dd if=/dev/urandom bs=1K count="$$size_kb" 2>/dev/null | tr -dc 'a-zA-Z0-9' > "bin/files/testFile_$${i}_$${size_kb}KB.txt"; \
	done

.PHONY: server
server: 
	./bin/server -config ./config/local.yml

.PHONY: client
client:
	./bin/client -config ./config/local.yml -path ./bin/files/testFile_1* -dir ./bin/files/
