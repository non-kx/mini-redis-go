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

test-html:
    go test -coverprofile coverage-html.out \
        ./internal/db/... \
        ./internal/network/... \
        ./internal/payload/... \
        ./internal/service/... \
        ./internal/tools/... \
        ./pkg/... \
	go tool cover -html=coverage-html.out

test-integration:
	go test -v ./test/...