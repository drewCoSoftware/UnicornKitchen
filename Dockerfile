FROM golang:1.16.3-buster

WORKDIR /go/src/app
COPY . .

RUN go get -d -v
RUN go install -v

CMD ["UnicornKitchen"]