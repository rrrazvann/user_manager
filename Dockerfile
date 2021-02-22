FROM golang:1.15.8-alpine

COPY ./src/ /go/src/

WORKDIR /go/src/user_manager/src

RUN go build -o main .

CMD ["./main"]