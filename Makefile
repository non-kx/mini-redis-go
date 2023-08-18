build:
	mkdir -p build
	cd cmd
	go build -o ../build ./...
	cd ..

test-integration:
	go test -v ./test/...