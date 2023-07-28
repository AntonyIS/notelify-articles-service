build:
	go build -o bin/notlify-user-svc

serve-prod: build
	./bin/notlify-user-svc -env=prod

serve-dev: build
	./bin/notlify-user-svc -env=dev