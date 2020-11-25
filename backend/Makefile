.PHONY: test
test:
	go test -race -v -coverprofile coverage.out ./...

.PHONY: build
build:
	mkdir -p build
	go build -o ./build ./cmd/mf

.PHONY: clean
clean:
	rm -rf build
