build:
	go build -o bin/notelify-articles-service
	
serve: build
	./bin/notelify-articles-service

test:
	Env=test go test -v ./...


