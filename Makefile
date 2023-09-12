build:
	go build -o bin/notelify-articles-service

serve-prod: build
	./bin/notelify-articles-service -env=prod

serve-dev: build
	./bin/notelify-articles-service -env=dev