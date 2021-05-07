FROM golang:latest

WORKDIR /go/src/app
COPY . .

ENV GOPROXY=""

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["watch"]