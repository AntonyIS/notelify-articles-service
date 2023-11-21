# Notelify articles service


## Desciption 
Article service is one of many service in the notelify application. It has all business logic for the articles logic. Exposes APIs that get consumed by other services
Build using Go programming language, Postgres, Redis for data storage.
The Article service makes use of the Hexagonal architecture for better code testing, easy code modification.
├── notelify-articles-service
│   ├── config
│   │   └── config.go
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── adapters
│   │   │   ├── app
│   │   │   │   ├── controllers.go
│   │   │   │   ├── controllers_test.go
│   │   │   │   └── handler.go
│   │   │   ├── logger
│   │   │   │   └── standardLogger.go
│   │   │   └── repository
│   │   │       ├── dynamodb
│   │   │       │   └── dynamodb.go
│   │   │       ├── postgres
│   │   │       │   └── postgres.go
│   │   │       └── s3
│   │   │           └── s3.go
│   │   └── core
│   │       ├── domain
│   │       │   └── domain.go
│   │       ├── ports
│   │       │   └── ports.go
│   │       └── services
│   │           ├── services.go
│   │           └── service_test.go
│   ├── LICENSE
│   ├── main.go
│   ├── Makefile
│   └── README.md


## Table of content
- [Installation](#installation)
<!-- - [Usage](#usage)
- [Configuration](#configuration)
- [Contribution](#contribution)
- [License](#license)
- [Acknowledgements](#acknowledgements) -->


## 1.0 Installation 
### Clone the repository in your workspace
* git clone https://github.com/AntonyIS/notelify-articles-service.git
### Change directory into the application code
* cd notelify-articles-service
### Install Go depencies for the Article service
* go mod tidy
### Run the application
* make serve-dev
### Access the application
* Open your browser or API client like Postman and navigate to http://localhost:8001 to access the appliction API end points

## 2.0 Installation - Using Docker
### Clone the repository in your workspace
* git clone https://github.com/AntonyIS/notelify-articles-service.git
### Change directory into the application code
* cd notelify-articles-service
### Build Docker image 
* docker build -t notelify-article-service .
### Run the docker container 
* docker run -p 8001:8001 -d notelify-article-service
### Access the application
* Open your browser or API client like Postman and navigate to http://localhost:8001 to access the appliction API end points
### Stopping and removing the docker container
* docker stop $(docker ps -aq --filter ancestor=notelify-article-service)
* docker rm $(docker ps -aq --filter ancestor=notelify-article-service)