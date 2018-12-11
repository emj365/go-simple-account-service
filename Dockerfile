FROM golang:1.11
WORKDIR /go/src/app
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
COPY ./Gopkg.lock Gopkg.toml /go/src/app/
RUN dep ensure -vendor-only
COPY . /go/src/app
CMD ["go", "build", "main.go"]

FROM alpine:latest
EXPOSE 8000
WORKDIR /root/
COPY --from=0 /go/src/app/app .
CMD ["./app"]
