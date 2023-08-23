build:
	mkdir -p build
	cd cmd
	go build -o ../build ./...
	cd ..

test-coverage:
	gocov test \
        ./internal/db/... \
        ./internal/network/... \
		./internal/payload/... \
        ./internal/service/... \
		./internal/tools/... \
		./pkg/... \
	| gocov report

test-integration:
	go test -v ./test/...