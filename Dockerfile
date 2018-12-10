FROM golang:1.11

EXPOSE 8000

WORKDIR /go/src/github.com/emj365/account

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
COPY ./Gopkg.lock Gopkg.toml /go/src/github.com/emj365/account/
RUN dep ensure -vendor-only

COPY . /go/src/github.com/emj365/account

CMD go run main.go
