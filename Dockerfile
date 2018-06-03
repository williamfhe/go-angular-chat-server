FROM golang:latest

WORKDIR /go/src/github.com/williamfhe/go-angular-chat-server 
ADD . .

RUN go get -v ./
RUN go build .

EXPOSE 8080

CMD ["go-angular-chat-server"]
