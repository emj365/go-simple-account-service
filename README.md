# go-simple-account-service

Simple Account Service in Golang

## Go Doc

https://godoc.org/github.com/emj365/go-simple-account-service

## Get Source

```bash
go get github.com/emj365/go-simple-account-service
```

## Install Dependencies

### MacOS

```bash
brew install dep
dep ensure
```

### Ohter

```bash
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
dep ensure
```

## Run

```bash
go run main.go
```

## Run with docker

```bash
cp .env.sample .env
docker build -t go-simple-account-service:latest .
docker run \
       -p 8000:8000 \
       -v $(pwd)/.env:/root/.env \
       go-simple-account-service:latest
```

## Request with curl

```bash
→ curl localhost:8000/users -v -d '{"name":"mike", "password":"passwd"}'
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8000 (#0)
> POST /users HTTP/1.1
> Host: localhost:8000
> User-Agent: curl/7.54.0
> Accept: */*
> Content-Length: 36
> Content-Type: application/x-www-form-urlencoded
>
* upload completely sent off: 36 out of 36 bytes
< HTTP/1.1 201 Created
< Content-Type: application/json
< Date: Tue, 11 Dec 2018 10:46:46 GMT
< Content-Length: 159
<
{"id":18,"created_at":"2018-12-11T18:46:46.482784+08:00","updated_at":"2018-12-11T18:46:46.482784+08:00","deleted_at":null,"name":"mike","password":"*******"}
* Connection #0 to host localhost left intact
```

```bash
→ curl localhost:8000/auth -v -d '{"name":"mike", "password":"passwd"}'
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8000 (#0)
> POST /auth HTTP/1.1
> Host: localhost:8000
> User-Agent: curl/7.54.0
> Accept: */*
> Content-Length: 36
> Content-Type: application/x-www-form-urlencoded
>
* upload completely sent off: 36 out of 36 bytes
< HTTP/1.1 200 OK
< Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NDQ1Mjg4NzAsInN1YiI6IjE4In0.ETwBxo08XiC6ygMN_OrxiCXd79sxIevwXSYwVxMpqc8
< Content-Type: application/json
< Date: Tue, 11 Dec 2018 10:47:50 GMT
< Content-Length: 131
<
{"jwt":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NDQ1Mjg4NzAsInN1YiI6IjE4In0.ETwBxo08XiC6ygMN_OrxiCXd79sxIevwXSYwVxMpqc8"}
* Connection #0 to host localhost left intact
```
