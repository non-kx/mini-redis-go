build:
	mkdir -p build
	cd cmd
	go build -o ../build ./...
	cd ..