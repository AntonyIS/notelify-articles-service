build:
	go build -o bin/notelify-articles-service
	
serve: build
	./bin/notelify-articles-service

test:
	go test -v -tags=myenv ./...
	Env=dev go test -v -tags=myenv ./...


