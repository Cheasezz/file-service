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

.PHONY: gen-test-file
gen-test-file:
	cat /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c 100000 > bin/testFile.txt

.PHONY: server
server: 
	./bin/server -config ./config/local.yml

.PHONY: client
client:
	./bin/client -config ./config/local.yml -path ./bin/testFile.txt
