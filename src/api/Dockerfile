FROM golang:alpine

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy

RUN go mod vendor

RUN go build -o api_main

EXPOSE 8080

CMD ["./api_main"]
