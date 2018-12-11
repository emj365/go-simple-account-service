FROM golang:1.11
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
WORKDIR /go/src/github.com/emj365/go-simple-account-service
COPY ./Gopkg.lock Gopkg.toml ./
RUN dep ensure -vendor-only
COPY . ./
# RUN CGO_ENABLED=0 go build -o main // panic: failed to connect database
RUN go build -o main

FROM alpine:latest
RUN apk add --no-cache libc6-compat
EXPOSE 8000
WORKDIR /root/
COPY --from=0 /go/src/github.com/emj365/go-simple-account-service/main .
CMD ["./main"]
