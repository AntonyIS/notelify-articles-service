FROM golang:alpine as builder

ENV GO111MODULE=on

RUN apk update && apk add --no-cache git

WORKDIR /app

ARG ENV

ENV ENV=${ENV}

COPY go.mod go.sum  ./

RUN go mod download 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./src .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/src .

EXPOSE 8001

CMD ["./src"]
