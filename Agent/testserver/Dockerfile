FROM golang:alpine

WORKDIR /usr/src/app

COPY . .

RUN go mod vendor

RUN go mod tidy

RUN go build -o testserver

EXPOSE 4201

CMD ["./testserver"]