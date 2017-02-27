FROM golang:1.8-alpine
RUN mkdir /app
ENV GOPATH=/app
RUN apk update && apk add --no-cache git bash openssh && \
go get github.com/el-komandante/gochat && \
cd $GOPATH && \
cd src/github.com/el-komandante/gochat && \
go build main.go && \
cp main /app && \
rm main
EXPOSE 8000
CMD ["/app/main"]
