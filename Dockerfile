FROM golang:1.8-alpine
RUN apk update && apk add --no-cache git && \
go get github.com/el-komandante/gochat && \
cd $GOPATH && \
cd src/github.com/el-komandante/gochat && \
go run main.go
